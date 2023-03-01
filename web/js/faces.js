const faceWrapper = document.getElementById("face");

let face;
const params = new URL(location.href).searchParams;
const avatar = getCookie('avatar');
const lang = params.get('lang');
const lang2 = navigator.language
var voices = window.speechSynthesis.getVoices();
var voice = '';

/*
 * Check for browser support
 */
var supportMsg = document.getElementById('errlabel');

if ('speechSynthesis' in window) {
	supportMsg.innerHTML = 'Your browser <strong>supports</strong> speech synthesis.';
} else {
	supportMsg.innerHTML = 'Sorry your browser <strong>does not support</strong> speech synthesis.<br>Try this in <a href="https://www.google.co.uk/intl/en/chrome/browser/canary.html">Chrome Canary</a>.';
	supportMsg.classList.add('not-supported');
}


// Fetch the list of voices and populate the voice options.
function loadVoices() {
  // Fetch the available voices.
	
  var voices = window.speechSynthesis.getVoices();
}

// Chrome loads voices asynchronously.
window.speechSynthesis.onvoiceschanged = function(e) {
  loadVoices();
};



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

var voicesApp = {"-":0};

function GetSpeechSynthesisId(voice) {
  console.log(voicesApp);
  for (var key in voicesApp) {
    if (key == voice || key == lang) {
      return voicesApp[key]
    }
  }
  for (var i=0;i < voices.length;i++) {
    if (voices[i].name.includes(voice)) {
      voicesApp[voice] = i;
      return i
    }
  }
  for (var i=0;i < voices.length;i++) {
    console.log(voices[i].lang, lang, i);
    if (voices[i].lang.startsWith(lang)) {
      voicesApp[lang] = i;
      return i
    }
  }
  return 0
}

function Speak(text) {
//const Speak = async(text) => {
  let Speech = new SpeechSynthesisUtterance();
  Speech.addEventListener('start', handleStartSpeechEvent);
  Speech.addEventListener('end', handleEndSpeechEvent);
  id =  GetSpeechSynthesisId(voice);
  Speech.voice = voices[id];
  Speech.text = text;
  console.log(voice, id);
  window.speechSynthesis.speak(Speech);
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

