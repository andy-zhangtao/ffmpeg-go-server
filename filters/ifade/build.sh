#!/bin/bash
ln -s /root/ffmpeg_build/include/* /usr/local/include && \
ln -s /root/ffmpeg_build/lib/* /usr/local/lib && \
cmake . && \
make