const faceWrapper = document.getElementById("face");

let face;
const params = new URL(location.href).searchParams;
const avatar = params.get('avatar');
const lang = params.get('lang');
const lang2 = navigator.language
const voices = window.speechSynthesis.getVoices();

var SpeechSynthesisId = 0;

(async() => {

  const getVoices = (voiceName = "") => {
    return new Promise(resolve => {
      window.speechSynthesis.onvoiceschanged = e => {
        resolve(window.speechSynthesis.getVoices());
      }
      window.speechSynthesis.getVoices();
    })
  }
  const voices = await getVoices();
  console.log(voices);
})();

if (avatar.length <= 1) {
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

const updateDisplay  = () => {  
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



const Speak = async(text) => {
  // Testing for browser support
	var speechSynthesisSupported = 'speechSynthesis' in window;
  let Speech = new SpeechSynthesisUtterance();
  console.log(SpeechSynthesisId);
  Speech.addEventListener('start', handleStartSpeechEvent);
  Speech.addEventListener('end', handleEndSpeechEvent);
  Speech.voice = speechSynthesis.getVoices()[ SpeechSynthesisId];
  Speech.text = text;
  speechSynthesis.speak(Speech);  
}
var speakingMode = false;

const handleStartSpeechEvent = async() => {
  speakingMode = true;
  document.body.style.cursor = 'wait';
  while (speakingMode) {
    await sleep(100);
    speaking();
  }
}


const handleEndSpeechEvent = async() => {
  speakingMode = false;
  document.body.style.cursor = 'default';
}

updateDisplay();

