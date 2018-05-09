package docker

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"

	"bufio"
	"encoding/json"
	"strings"

	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
)

// 创建镜像所使用的结构体
type CreateContainerConf struct {
	Host        string //连接远程客户端的信息 eg：tcp://222.24.63.117:9000
	ServiceName string //服务名称，用作最后镜像名称
	Image       string //镜像名称
	//config所需参数
	WorkDir  string   //cmd执行的工作目录
	HostName string   //容器的hostname
	Env      []string //环境变量
	Cmd      []string //cmd内容
	//hostconfig所需参数
	HostList []string //host:ip的格式，存入host对应ip的解析
	Dns      []string //dns服务器ip地址
	Volume   []string //path:hostpath，容器内部目录和宿主机目录的挂载
	//config和hostconfig综合参数
	Expose []string //port:hostport，容器端口和宿主机端口的映射
}

// 写入日志所用的json格式
type LogJson struct {
	machineId  int64  `json:"machine_id"`
	machineLog string `json:"machine_log"`
}

// 创建远程连接的客户端
func newClient(host string) (*client.Client, error) {
	// cli, err := client.NewClientWithOpts()
	cli, err := client.NewClient(host, "1.24", nil, nil)
	if err != nil {
		return nil, err
	}
	return cli, err
}

// 创建容器
func CreateContainer(param CreateContainerConf) (string, error) {
	cli, err := newClient(param.Host)
	if err != nil {
		return "", err
	}

	var exposePorts nat.PortSet
	exposePorts = make(map[nat.Port]struct{})
	var portBindings nat.PortMap
	portBindings = make(map[nat.Port][]nat.PortBinding)
	//获取config中的ExposedPorts和hostconfig中的PortBindings
	for _, expose := range param.Expose {
		str := strings.Split(expose, ":")
		if len(str) != 2 {
			log.Errorln("expose param error: ", expose)
			return "", errors.Errorf("expose param error: %s", expose)
		}

		str[1] = strings.Replace(str[1], " ", "", -1)
		str[1] += "/tcp"
		containerPort := nat.Port(str[1])
		exposePorts[containerPort] = struct{}{}

		portBind := []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: str[0],
			},
		}
		portBindings[containerPort] = portBind
	}
	//fmt.Println(exposePorts, "\n", portBindings)

	ctx := context.Background()
	containerBody, err := cli.ContainerCreate(ctx,
		&container.Config{
			User:         "root",
			Image:        param.Image,
			WorkingDir:   param.WorkDir,
			Hostname:     param.HostName,
			Env:          param.Env,
			Cmd:          param.Cmd,
			ExposedPorts: exposePorts,
		}, &container.HostConfig{
			Resources: container.Resources{
				NanoCPUs: 2,
				Memory:   512000000,
			},
			ExtraHosts:      param.HostList,
			DNS:             param.Dns,
			PortBindings:    portBindings,
			Binds:           param.Volume,
			PublishAllPorts: true,
		}, nil, param.ServiceName)

	if err != nil {
		log.WithField("err", err.Error()).Error("docker container create failure")
		return "", err
	}
	return containerBody.ID, nil
}

// 运行容器
func StartContainer(host string, containerId string) error {
	cli, err := newClient(host)
	if err != nil {
		return err
	}
	ctx := context.Background()

	err = cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
	return err
}

// 停止容器
func StopContainer(host string, containerID string) error {
	cli, err := newClient(host)
	if err != nil {
		return err
	}
	timeout := time.Second * 20
	err = cli.ContainerStop(context.Background(), containerID, &timeout)
	return err
}

// 移除容器
func RemoveContainer(host string, containerID string, force bool, removeVolumes bool, removeLinks bool) error {
	cli, err := newClient(host)
	if err != nil {
		return err
	}
	ctx := context.Background()

	options := types.ContainerRemoveOptions{Force: force, RemoveVolumes: removeVolumes, RemoveLinks: removeLinks}
	if err := cli.ContainerRemove(ctx, containerID, options); err != nil {
		return err
	}
	return nil
}

// 杀死容器
func KillContainer(host string, containerId string) error {
	cli, err := newClient(host)
	if err != nil {
		return err
	}
	ctx := context.Background()

	err = cli.ContainerKill(ctx, containerId, "SIGKILL")
	return err
}

// 显示容器列表
func ListContainers(host string) ([]types.Container, error) {
	cli, err := newClient(host)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// 显示镜像列表
func ListImages(host string) ([]types.ImageSummary, error) {
	cli, err := newClient(host)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, nil
	}

	return images, nil
}

// 获取pull 镜像时的json串格式
type OutJson struct {
	Status string `json:"status"`
	Id     string `json:"id"`
}

// 每一行日志信息的格式
type RowJson struct {
	JsonStr OutJson
	Mess    string
}

func PullImage(host string, imageName string) (string, error) {
	cli, err := newClient(host)
	if err != nil {
		return "", err
	}
	ctx := context.Background()

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	defer out.Close()

	list := make([]RowJson, 0)
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		//fmt.Printf("[%s]\n", scanner.Text())
		data := scanner.Text()
		var jsonStr OutJson

		if err := json.Unmarshal([]byte(data), &jsonStr); err != nil {
			panic(err)
		}
		// 说明是开头或结尾
		if jsonStr.Id == "" {
			list = append(list, RowJson{JsonStr: jsonStr, Mess: data})
		} else {
			flag := -1
			for index, row := range list {
				if row.JsonStr.Id == jsonStr.Id {
					flag = index
				}
			}
			if flag == -1 {
				list = append(list, RowJson{JsonStr: jsonStr, Mess: data})
			} else {
				newList := list[:flag]
				newList = append(newList, RowJson{JsonStr: jsonStr, Mess: data})
				newList = append(newList, list[flag+1:]...)
				list = newList
			}
		}
		// 这里是写入数据库的内容
		//fmt.Println("=======上")
		//for _, row := range list {
		//	fmt.Println(row.Mess)
		//}
		//fmt.Println("=======下")
	}
	var mess string
	for _, row := range list {
		mess += row.Mess
		mess += "<br>"
	}

	return mess, nil
}

// 打印镜像日志
func PrintLogContainer(host string, containerID string) error {
	cli, err := newClient(host)
	if err != nil {
		return err
	}
	ctx := context.Background()

	options := types.ContainerLogsOptions{ShowStdout: true}
	out, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	return nil
}
