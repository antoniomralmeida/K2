const faceWrapper = document.getElementById("face");

let face;
const params = new URL(location.href).searchParams;
const avatar = params.get('avatar');
const lang = params.get('lang');
const lang2 = navigator.language
var voices = window.speechSynthesis.getVoices();
var voice = '';

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

var voicesApp = {"Mark":1};

function GetSpeechSynthesisId(voice) {
  console.log(voicesApp);
  for (var key in voicesApp) {
    if (key == voice) {
      return voicesApp[key]
    }
  }
  for (var i=0;i < voices.length;i++) {
    if (voices[i].name.includes(voice)) {
      voicesApp[voice] = i;
      return i
    }
  }
  return 1
}

const Speak = async(text) => {
  // Testing for browser support
	var speechSynthesisSupported = 'speechSynthesis' in window;
  let Speech = new SpeechSynthesisUtterance();
  while (voices.length==0) { 
    voices = window.speechSynthesis.getVoices();
    await sleep(50);
  }
  Speech.addEventListener('start', handleStartSpeechEvent);
  Speech.addEventListener('end', handleEndSpeechEvent);
  id =  GetSpeechSynthesisId(voice);
  Speech.voice = voices[id];
  Speech.text = text;
  console.log(voice, id);
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
