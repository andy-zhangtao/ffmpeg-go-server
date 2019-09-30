cmake_minimum_required(VERSION 3.10)
project(ifade C)

set(CMAKE_C_STANDARD 99)
set(CMAKE_C_FLAGS -pthread)
#set(CMAKE_PREFIX_PATH /root/ffmpeg_build)

find_path(AVCODEC_INCLUDE_DIR libavcodec/avcodec.h DOC "find avcodec header file")
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

add_library(ifade copy.c copy.h)

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