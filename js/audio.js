const recordAudio = () =>
  new Promise(async resolve => {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    const mediaRecorder = new MediaRecorder(stream);
    const audioChunks = [];

    mediaRecorder.addEventListener("dataavailable", event => {
      audioChunks.push(event.data);
    });

    const start = () => mediaRecorder.start();
    const reset = () => audioChunks.length = 0;

    const stop = () =>
      new Promise(resolve => {
        mediaRecorder.addEventListener("stop", () => {
          const audioBlob = new Blob(audioChunks);
          const audioUrl = URL.createObjectURL(audioBlob);
          const audio = new Audio(audioUrl);
          const play = () => audio.play();
          resolve({ audioBlob, audioUrl, play });
        });
        console.log("audio chunks are:");
        console.log(audioChunks);
        mediaRecorder.stop();
      });

    resolve({ start, stop, reset});
    
  });

const sleep = time => new Promise(resolve => setTimeout(resolve, time));

(async () => {
  const recorder = await recordAudio();

  document.getElementById('toggleRecord').onclick = function(){
    if(document.getElementById('toggleRecord').innerText == "Start Recording"){
        console.log("start recording");  
        recorder.start();
        document.getElementById('toggleRecord').innerText = "Stop Recording";
        document.getElementById("resetRecording").disabled = true;
    }else{

        (async () => {
            const audio = await recorder.stop();
            console.log("stop recording");
            audio.play();
            document.getElementById('toggleRecord').innerText = "Start Recording";
            document.getElementById("resetRecording").disabled = false;
        })();  
    }
  }

  document.getElementById('resetRecording').onclick = function(){
    recorder.reset();
    document.getElementById("resetRecording").disabled = true;
  }
  
})();