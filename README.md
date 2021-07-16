
# ServerManagerSystem
一个基于Golang开发的远程服务管理系统，使用RPC远程调度，管理服务器进程保持、数据库远程调度等

#### 说明
- RPC部分基于Hprose的RPC框架实现
- 本人第一次开源Golang项目，处于学习和探索阶段，尚有很多不完善
- 本项目主要功能是：
   1. 实现客户端简单配置自动RPC链接
   1. 实现客户端FRP NPS简单远程配置即可自动上线
   1. 实现远程shell/sftp的管理
   1. 实现远程和本地的进程保持
   1. 实现远程的SQL调度
   1. 简单的防火墙开关
   1. 基于VUE的管理页面
   1. 默认用户名：admin 密码：123456

#### 使用说明
- 以debian10为例安装golang

   ```shelL
   wget https://golang.org/dl/go1.16.6.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.16.6.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin'>> 〜/ .profile
   source ~/.profile
   go version
   ```

- 编译并使用本程序
    ```shelL
      git clone https://github.com/uouuou/ServerManagerSystem.git
      sh build.sh
      cp -r ./web ./output
      cd output
      chemo -R 755 ServerManagerSystem_linux_amd64_upx 
      ./ServerManagerSystem_linux_amd64_upx  install
      service serverManager start
      journalctl -u ServerManager.service -f
    ```

#### 界面展示

![](https://i.lioil.cc/o0o/2021/07/16/cb05fdb0a62eab34.png)

![](https://i.lioil.cc/o0o/2021/07/16/44a0afb55adcce97.png)

![](https://i.lioil.cc/o0o/2021/07/16/03e8c5f739f7037e.png)
