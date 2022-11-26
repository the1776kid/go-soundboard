#!/bin/bash
targetSampleRate=48000
audioDir=$1
# TODO is $1 is nil print help
echo "${audioDir}"
for i in $(ls ${audioDir}) ; do
  sr=$(ffprobe -hide_banner -loglevel error -show_entries stream=sample_rate ${audioDir}${i} | awk 'FNR == 2 {print}')
  echo "${audioDir}${i} ${sr}"
  if [[ ${sr} != *"${targetSampleRate}"* ]];then
    echo "transcode ${i}, set sample rate ${targetSampleRate}hz"
    mv ${audioDir}${i} ${audioDir}tempFile.mp3
    ffmpeg -i ${audioDir}tempFile.mp3 -hide_banner -loglevel error -af aresample=resampler=soxr -ar ${targetSampleRate} ${audioDir}${i}
    rm ${audioDir}tempFile.mp3
  fi
done
