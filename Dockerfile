#添加JAVA启动的必要镜像
FROM java:8


#【配置时区统一】
RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

#设置开放端口号

#添加jar包，存放路径以及重命名
ADD  ./main  //main


#添加进入docker容器后的目录
WORKDIR   /

#修改文件的创建修改时间
RUN bash -c 'touch /main'
RUN bash -c 'chmod 777 /main'
#启动容器执行命令。 -Djava.security.egd=file:/dev/./urandom 可以缩短tomcat启动时间
ENTRYPOINT ["./main"]
