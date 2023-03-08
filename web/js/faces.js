const faceWrapper = document.getElementById("face");

let face;
const params = new URL(location.href).searchParams;
const avatar = getCookie('avatar');
const lang = params.get('lang');
var arrayVoices = []
Synthesis.getVoices();
var voice = '';

/*
 * Check for browser support
 */
var supportMsg = document.getElementById('errlabel');

if (!'speechSynthesis' in window)  {
	supportMsg.innerHTML = 'Sorry your browser <strong>does not support</strong> speech synthesis.<br>Try this in <a href="https://www.google.co.uk/intl/en/chrome/browser/canary.html">Chrome Canary</a>.';
}


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

var voicesApp = { "-": -1 };

function GetSpeechSynthesisId(voice) {

  console.log(voicesApp, voice, lang, arrayVoices.length);

  for (var key in voicesApp) {
    if (key == voice || key == lang) {
      return voicesApp[key];
    }
  }

  for (var i = 0; i < arrayVoices.length; i++) {
    if (arrayVoices[i].name.includes(voice)) {
      voicesApp[voice] = i;
      return i;
    }
  }
  for (var i = 0; i < arrayVoices.length; i++) {
    console.log(arrayVoices[i].lang, lang, i);
    if (arrayVoices[i].lang.startsWith(lang)) {
      voicesApp[lang] = i;
      return i;
    }
  }
  return -1;
}

const Speak = async (text) => {
  // Testing for browser support
  var speechSynthesisSupported = 'speechSynthesis' in window;
  let Speech = new SpeechSynthesisUtterance();
  while (arrayVoices.length==0) { 
    arrayVoices = window.speechSynthesis.getVoices();
    await sleep(50);
  }
  Speech.addEventListener('start', handleStartSpeechEvent);
  Speech.addEventListener('end', handleEndSpeechEvent);
  id = GetSpeechSynthesisId(voice);
  Speech.voice = arrayVoices[id];
  Speech.text = text;
  console.log(voice, id);
  if (id >= 0) {
    speechSynthesis.speak(Speech);
  }
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
