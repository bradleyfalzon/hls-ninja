#!/bin/bash -eu

DURATION=300
MASTER="testsrc"

function run() {
    mkdir -p ${MASTER}/${RES} || true
    ffmpeg -re -f lavfi -i testsrc=size=${RES}:rate=25 -f lavfi -i "sine=frequency=1000"  -af volume=0.1 -vf "drawtext=fontfile=/usr/share/fonts/dejavu/DejaVuSansMono.ttf:text='Seconds\: %{pts} Res\:${RES}':fontsize=30:fontcolor="white":boxcolor=0x00000088:box=1:x=(w-tw)/2:y=h*0.27:boxborderw=20:bordercolor=0x00000088" -s ${RES} -b:v ${RATE} -maxrate ${RATE} -bufsize ${RATE} -c:v libx264 -pix_fmt yuv420p -profile:v main -level 3.1 -bf 0 -g 50 -keyint_min 50 -sc_threshold 0 -c:a aac -ac 2 -ar 48000 -strict -2 -ab 64k -f hls -hls_time 6  -hls_segment_filename "${MASTER}/${RES}/%02d.ts" -hls_base_url "${RES}/" -hls_list_size 0 -t ${DURATION} ${MASTER}/${RES}.m3u8
    echo "#EXT-X-STREAM-INF:BANDWIDTH=${RATE},RESOLUTION=${RES}" >> ${MASTER}.m3u8
    echo "${MASTER}/${RES}.m3u8" >> ${MASTER}.m3u8
}

function makeMaster() {
    echo "#EXTM3U" > ${MASTER}.m3u8
    echo "#EXT-X-VERSION:3" >> ${MASTER}.m3u8
}

makeMaster

RES=1920x1080
RATE=1200000
run

RES=1280x720
RATE=900000
run

RES=960x540
RATE=600000
run

RES=640x360
RATE=300000
run
