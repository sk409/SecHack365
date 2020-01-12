package main

import (
	"fmt"
	"os"
	"os/exec"
)

type dockerfile struct {
	text string
}

func newDockerfile(image string, username string) *dockerfile {
	f := dockerfile{}
	f.text += "FROM " + image + "\n"
	f.text += "RUN useradd " + username + "\n"
	f.text += `RUN yum -y install epel-release \
    && rpm -Uvh http://rpms.famillecollet.com/enterprise/remi-release-7.rpm \
    && echo -e "[epel]\nname=Extra Packages for Enterprise Linux 7 - \$basearch\n#baseurl=http://download.fedoraproject.org/pub/epel/7/\$basearch\nmirrorlist=https://mirrors.fedoraproject.org/metalink?repo=epel-7&arch=\$basearch\nfailovermethod=priority\nenabled=0\ngpgcheck=1\ngpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7" > etc/yum.repos.d/epel.repo \
    && yum -y install wget \
    && yum -y install vim \
    && yum -y install sudo \
    && wget -qO- https://github.com/yudai/gotty/releases/download/v0.0.12/gotty_linux_amd64.tar.gz | tar zx -C /usr/local/bin/ \
    && yum -y remove wget \
`
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

func (d *docker) runContainer(containername, imagename string, ports ...string) ([]byte, error) {
	args := []string{"container", "run", "-itd", "--name", containername}
	for _, port := range ports {
		args = append(args, "-p")
		args = append(args, port)
	}
	args = append(args, imagename)
	command := exec.Command("docker", args...)
	return command.Output()
}
