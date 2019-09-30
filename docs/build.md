# 构建
> 如何构建C编写的Filter

为了方便调用`AVFilter`， `FFmpeg-go-server`会通过`cgo`链接构建出的静态库，生成可执行文件。 所以就涉及到两部分:

+ 构建静态库
+ cgo链接静态库

### 构建静态库

#### 工程说明

编写filter complex完函数之后，应该会有两个文件:`xxx.c`和`xxx.h`。 

为了尽量减少`golang`和`C`之间数据转换，因此约定C暴露给golang的函数只需要提供两个参数即可，例如:

```C
    int copy(char *, char *);
```

在`golang`中可以通过下面的代码直接调用:
```cgo
    input_file := "/60.mp4"
    output_file := "/fade.mp4"

    ret := C.copy(C.CString(input_file), C.CString(output_file))
    if ret < 0 {
        fmt.Println("Fade Error. " + ret)
    }
```

除了`.c`和`.h`之外，还需要提供`CMakeLists.txt`用来生成`Makefile`。 `CMakeList.txt`可以直接使用下面的模板:

```makefile
cmake_minimum_required(VERSION 3.10)
project(<library name> C)

set(CMAKE_C_STANDARD 99)
set(CMAKE_C_FLAGS -pthread)

find_path(AVCODEC_INCLUDE_DIR libavcodec/avcodec.h)
find_library(AVCODEC_LIBRARY avcodec)

find_path(AVFILTER_INCLUDE_DIR libavfilter/avfilter.h)
find_library(AVFILTER_LIBRARY avfilter)

find_path(AVFORMAT_INCLUDE_DIR libavformat/avformat.h)
find_library(AVFORMAT_LIBRARY avformat)

find_path(AVUTIL_INCLUDE_DIR libavutil/avutil.h)
find_library(AVUTIL_LIBRARY avutil)

find_path(AVDEVICE_INCLUDE_DIR libavdevice/avdevice.h)
find_library(AVDEVICE_LIBRARY avdevice)

find_path(SWSCALE_INCLUDE_DIR libswscale/swscale.h)
find_library(SWSCALE_LIBRARY swscale)

add_library(<library name> xxx.c xxxx.h)

INCLUDE_DIRECTORIES(
        /usr/local/include

)

LINK_DIRECTORIES(
        /usr/local/lib
        /lib
)

target_include_directories(ifade PRIVATE ${AVCODEC_INCLUDE_DIR} ${AVFILTER_INCLUDE_DIR} ${AVFORMAT_INCLUDE_DIR} ${AVUTIL_INCLUDE_DIR} ${AVDEVICE_INCLUDE_DIR} ${SWSCALE_INCLUDE_DIR})
target_link_libraries(ifade
        PRIVATE
        ${AVCODEC_LIBRARY}
        ${AVFILTER_LIBRARY}
        ${AVFORMAT_LIBRARY}
        ${AVUTIL_LIBRARY}
        ${AVDEVICE_LIBRARY}
        ${SWSCALE_LIBRARY}
        )
```

> `add_library`表示构建出静态库。

最后需要提供构建脚本: `build.sh`. 这个脚本用来执行`cmake`和`make`，一般来说内容是下面的样子:

```
#!/bin/bash
ln -s /root/ffmpeg_build/include/* /usr/local/include && \
ln -s /root/ffmpeg_build/lib/* /usr/local/lib && \
cmake . && \
make
```

一个完整的`AVFilter`目录结构如下:

![](https://tva1.sinaimg.cn/large/006y8mN6ly1g7hhrqq8l7j3056039aa1.jpg)

#### 构建过程

当上面的内容都准备之后，需要修改根目录下面的`Makefile`。 在文件中添加新`AVFilter`的目标。 

构建`AVFilter`的静态库使用的是`vikings/ffmpeg-ubuntu`([镜像说明](http://dockerfile.docs.devexp.cn/ffmpeg-ubuntu.html)).

使用此镜像时，需要指定`workdir`(配合`build.sh`)。


### CGO

当成功构建静态库之后，就可以在`golang`中通过`cgo`调用静态库。 以下略过golang代码说明，只说明`cgo`构建参数。

以`ifade`为例， 如果采用上面的`CMake`参数，最终构建的静态库名称应该是`libifade.a`。 假设此静态库与`main.go`在同级目录，并且在`main.go`中调用函数，那么CGO参数应该是:
```cgo
// #cgo pkg-config: libavfilter libavdevice
// #cgo LDFLAGS: libifade.a
// #include "copy.h"
```

`pkg-config`会尝试自动展开指定静态库的`cflags`和`libs`参数，而因为`libifade.a`是我们自己构建出来的，并没有`.pc`文件，所以无法通过`pkg-config`自动展开。需要单独指定其位置。
如果构建过程中提示找不到某函数定义，那么需要根据实际情况修改`pkg-config`参数。