# btchain
A block chain for BT of JVV

BTChain Baseed on tendermint 0.27.0-9c236ffd

## 安装cmake[使用cleveldb必选]
```
wget https://cmake.org/files/v3.6/cmake-3.6.2.tar.gz
./boostrap
make && make install
```

## 编译snappy[使用cleveldb必选]
```
wget https://github.com/google/snappy/archive/1.1.7.tar.gz
tar zxf 1.1.7.tar.gz
cd snappy-1.1.7
编辑CMakeLists.txt将BUILD_SHARED_LIBS开关打开（默认生成静态库，编译cleavel要静态库）
mkdir build && cd build && cmake ../
make install
```

## 编译cleaveldb[使用cleveldb必选]
```
wget https://github.com/google/leveldb/archive/v1.20.tar.gz && \
  tar -zxvf v1.20.tar.gz && \
  cd leveldb-1.20/ && \
  make && \
  sudo cp -r out-static/lib* out-shared/lib* /usr/local/lib/ && \
  cd include/ && \
  sudo cp -r leveldb /usr/local/include/ && \
  sudo ldconfig && \
  rm -f v1.20.tar.gz
  ```

## 编译链程序
- 先按照tendermint标准编译流程进行编译，确保单独编译tendermint可以成功
- 修改tendermint/proxy/client.go,func:DefaultClientCreator,加入以下分支
```
	case "bt":
		return NewLocalClientCreator(btchain.NewBTApplication())
```
导入包
```
"github.com/axengine/btchain"
```
- 确保gcc可用，因为要开启cgo
- 进入tendermint目录执行`CGO_LDFLAGS="-lsnappy" make build_c`，生成的二进制文件在build目录
- ldd tendermint,查看依赖的动态库是否都能找到

## 使用
- 配置文件config.toml与可执行程序在同一目录
```
[genesis]
account = "0x061a060880BB4E5AD559350203d60a4349d3Ecd6"
amount = "10000000000"


[db]
type = "sqlite3"
path = "./data/"

[log]
env = "debug" # production
path = "./log/" #需先创建
```
- 要使用cleveldb，需要修改.tendermint/config/config.toml 的backend
- 日志目录，与可执行程序同级的log目录，需先创建

## 编译API
- 进入github.com/axengine/btchain/api 执行`go build`即可
- 配置文件，在可执行程序所在目录的config目录下
```
bind = ":10000"
rpc = "127.0.0.1:26657"
writable = true #false时只有查询API
isAdmin = true #true时有validator更新API

[log]
path = "./log/"
```