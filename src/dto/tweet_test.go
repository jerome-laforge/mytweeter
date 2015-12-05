package dto_test

import (
	"fmt"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
)

func TestListImages(t *testing.T) {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		t.Fatal(err)
	}
	for _, img := range imgs {
		t.Log("ID: ", img.ID)
		t.Log("RepoTags: ", img.RepoTags)
		t.Log("Created: ", img.Created)
		t.Log("Size: ", img.Size)
		t.Log("VirtualSize: ", img.VirtualSize)
		t.Log("ParentId: ", img.ParentID)
		t.Log("-----------------------------")
	}
}

func TestCreateCassandra(t *testing.T) {
	_, err := CreateCassandraAndStart("toto", "9043")
	if err != nil {
		t.Fatal(err)
	}

}

func CreateCassandraAndStart(name, port string) (*docker.Container, error) {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	hostConf := &docker.HostConfig{
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port("9042/tcp"): []docker.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: port,
			}},
		},
	}
	/*hostConf := &docker.HostConfig{
		PublishAllPorts: true,
	}*/

	opts := docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: "spotify/cassandra:latest",
		},
		HostConfig: hostConf,
	}

	container, err := client.CreateContainer(opts)
	if err != nil {
		return nil, err
	}

	return container, client.StartContainer(container.ID, hostConf)
}

func PullImages(repoTag string, tag string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return err
	}

	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return err
	}

	for _, img := range imgs {
		for i := range img.RepoTags {
			if img.RepoTags[i] == repoTag {
				return nil
			}
		}
	}

	searchImages, err := client.SearchImages(repoTag)
	if err != nil {
		return err
	}

	if len(searchImages) == 0 {
		return fmt.Errorf("Image not found for %s", repoTag)
	}

	/*if len(searchImages) > 1 {
		for _, each := range searchImages {
			fmt.Println(each.Name)
		}
		return fmt.Errorf("Too images found for %s", repoTag)
	}*/

	if tag == "" {
		tag = "latest"
	}

	pullOption := docker.PullImageOptions{
		Repository: repoTag,
		Tag:        tag,
	}

	auth := docker.AuthConfiguration{}
	err = client.PullImage(pullOption, auth)
	if err != nil {
		return err
	}

	return nil
}
