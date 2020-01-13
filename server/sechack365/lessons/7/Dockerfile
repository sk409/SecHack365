FROM centos:7
RUN useradd test
RUN yum -y install epel-release \
    && rpm -Uvh http://rpms.famillecollet.com/enterprise/remi-release-7.rpm \
    && echo -e "[epel]\nname=Extra Packages for Enterprise Linux 7 - \$basearch\n#baseurl=http://download.fedoraproject.org/pub/epel/7/\$basearch\nmirrorlist=https://mirrors.fedoraproject.org/metalink?repo=epel-7&arch=\$basearch\nfailovermethod=priority\nenabled=0\ngpgcheck=1\ngpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7" > etc/yum.repos.d/epel.repo \
    && yum -y install wget \
    && yum -y install vim \
    && yum -y install sudo \
    && wget -qO- https://github.com/yudai/gotty/releases/download/v0.0.12/gotty_linux_amd64.tar.gz | tar zx -C /usr/local/bin/ \
	&& yum -y remove wget \
	&& echo -e "preferences{\nbackground_color = \"rgb(255, 255, 255)\"\nforeground_color = \"rgb(16, 16, 16)\"}" >> ~/.gotty \
    && echo "test ALL=NOPASSWD: ALL" >> /etc/sudoers
USER test