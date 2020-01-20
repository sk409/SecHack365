package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type dockerfile struct {
	text string
}

func newDockerfile(image string, username string) *dockerfile {
	f := dockerfile{}
	f.text += "FROM " + image + "\n"
	f.text += "RUN useradd " + username + "\n"
	if strings.Contains(image, "centos") {
		f.text += `RUN yum -y install epel-release \
    && rpm -Uvh http://rpms.famillecollet.com/enterprise/remi-release-7.rpm \
    && echo -e "[epel]\nname=Extra Packages for Enterprise Linux 7 - \$basearch\n#baseurl=http://download.fedoraproject.org/pub/epel/7/\$basearch\nmirrorlist=https://mirrors.fedoraproject.org/metalink?repo=epel-7&arch=\$basearch\nfailovermethod=priority\nenabled=0\ngpgcheck=1\ngpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7" > etc/yum.repos.d/epel.repo \
    && yum -y install wget \
    && yum -y install vim \
    && yum -y install sudo \
`
		f.text += "    && wget -qO- https://github.com/yudai/gotty/releases/download/v0.0.12/gotty_linux_amd64.tar.gz | tar zx -C /usr/local/bin/ \\\n"
		f.text += fmt.Sprintf("    && echo -e \"preferences{\\nbackground_color = \\\"rgb(255, 255, 255)\\\"\\nforeground_color = \\\"rgb(16, 16, 16)\\\"\\n}\" >> home/%s/.gotty \\\n", username)
	} else if strings.Contains(image, "ubuntu") {
		f.text += `RUN apt-get update \
	&& apt-get -y install wget \
	&& apt-get -y install vim \
	&& apt-get -y install sudo \
`
		f.text += fmt.Sprintf("    && mkdir /home/%s \\\n", username)
		f.text += "    && wget -qO- https://github.com/yudai/gotty/releases/download/v0.0.12/gotty_linux_amd64.tar.gz | tar zx -C /usr/local/bin/ \\\n"
		f.text += fmt.Sprintf("    && echo \"preferences{\\nbackground_color = \\\"rgb(255, 255, 255)\\\"\\nforeground_color = \\\"rgb(16, 16, 16)\\\"\\n}\" >> home/%s/.gotty \\\n", username)
	}
	f.text += fmt.Sprintf("    && echo \"%s ALL=NOPASSWD: ALL\" >> /etc/sudoers\n", username)
	f.text += fmt.Sprintf("USER %s", username)
	return &f
}

func (d *dockerfile) write(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(d.text))
	return err
}

type docker struct {
}

func (d *docker) buildImage(imagename, directoryPath string) error {
	command := exec.Command("docker", "image", "build", "-t", imagename, directoryPath)
	err := command.Start()
	if err != nil {
		return err
	}
	return command.Wait()
}

func (d *docker) commit(id, imagename string) error {
	command := exec.Command("docker", "container", "commit", id, imagename)
	err := command.Start()
	if err != nil {
		return err
	}
	return command.Wait()
}

func (d *docker) exec(id string, args []string, command ...string) ([]byte, error) {
	args = append([]string{"container", "exec"}, args...)
	args = append(args, id)
	args = append(args, command...)
	cmd := exec.Command("docker", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		log.Println(stderr.String())
	}
	return output, err
}

func (d *docker) image(id string) ([]byte, error) {
	b, err := d.inspect(id, "--format", "{{.Image}}")
	if err != nil {
		return nil, err
	}
	c := bytes.Split(b, []byte(":"))
	return c[1], nil
}

func (d *docker) inspect(id string, args ...string) ([]byte, error) {
	args = append([]string{"container", "inspect", id}, args...)
	command := exec.Command("docker", args...)
	return command.Output()
}

func (d *docker) kill(id string) ([]byte, error) {
	command := exec.Command("docker", "container", "kill", id)
	return command.Output()
}

func (d *docker) port(id string) ([]byte, error) {
	command := exec.Command("docker", "container", "port", id)
	return command.Output()
}

func (d *docker) removeContainer(id string) ([]byte, error) {
	command := exec.Command("docker", "container", "rm", id)
	return command.Output()
}

func (d *docker) runContainer(containername, imagename string, ports ...string) ([]byte, error) {
	args := []string{"container", "run", "-itd", "--name", containername}
	for _, port := range ports {
		args = append(args, "-p")
		args = append(args, port)
	}
	args = append(args, imagename)
	command := exec.Command("docker", args...)
	command.Stderr = os.Stderr
	return command.Output()
	// var e = bytes.Buffer{}
	// command.Stderr = &e
	// err := command.Start()
	// if err != nil {
	// 	log.Println(e.String())
	// 	return nil, err
	// }
	// out, err := command.Output()
	// if err != nil {
	// 	log.Print(e.String())
	// 	return nil, err
	// }
	// log.Println(out)
	// return out, nil
}

func (d *docker) sendFile(id, src, dst string) error {
	command := exec.Command("docker", "cp", src, id+":"+dst)
	err := command.Start()
	if err != nil {
		return err
	}
	return command.Wait()
}
