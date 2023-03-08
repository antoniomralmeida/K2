const faceWrapper = document.getElementById("face");

let face;
const avatar = getCookie('avatar');


if (avatar == null) {
  face = faces.generate();
} else {
  try {
    face = JSON.parse(atob(avatar));
  } catch (error) {
    console.error(error);
    face = faces.generate();
  }
}


function sleep(milliseconds) {
  return new Promise(resolve => setTimeout(resolve, milliseconds));
}

const updateDisplay = () => {
  faces.display(faceWrapper, face);
};


const speaking = () => {
  if (face["mouth"].id == "smile") {
    face["mouth"].id = "mouth7";
  } else {
    face["mouth"].id = "smile";
  }
  updateDisplay();
}

const Speak = async(mp3File) => {
  console.log(mp3File);
  const audio = new Audio(mp3File);
  audio.onended = handleEndSpeechEvent;
  audio.onplay = handleStartSpeechEvent;
  audio.play();
}

var speakingMode = false;

const handleStartSpeechEvent = async () => {
  speakingMode = true;
  document.body.style.cursor = 'wait';
  while (speakingMode) {
    await sleep(100);
    speaking();
  }
}


const handleEndSpeechEvent = async () => {
  speakingMode = false;
  document.body.style.cursor = 'default';
}

updateDisplay();

