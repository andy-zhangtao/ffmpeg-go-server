ifade:
	docker run -it --rm -v $(PWD)/filters/ifade:/root/ifade -w /root/ifade vikings/ffmpeg-ubuntu sh build.sh

all: ifade