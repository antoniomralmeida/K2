const faceWrapper = document.getElementById("face");

let face;
const params = new URL(location.href).searchParams;
const avatar = params.get('avatar');
const lang = params.get('lang');
const lang2 = navigator.language
const voices = window.speechSynthesis.getVoices();

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



const TTS = async(text) => {
  // Testing for browser support
	var speechSynthesisSupported = 'speechSynthesis' in window;
  let msg = new SpeechSynthesisUtterance();
  i = 0;
  voices.forEach((voice) => {
    alert(i + ':'+voice);
    i=i+1;
  });
  msg.voice = speechSynthesis.getVoices()[1];
  if (lang == 'pt') {
    msg.voice = speechSynthesis.getVoices()[0];
  } else if (lang == null && lang2 == 'pt-BR') {
    msg.voice = speechSynthesis.getVoices()[0];
  }
  msg.text = text;
  speechSynthesis.speak(msg);  
}

const Speak = async (text) => {
  await sleep(50);
  document.body.style.cursor = 'wait';
  interval = text.length * 100;
  TTS(text);
  var d = new Date();
  var begin = d.getTime()
  while (true) {
    speaking();
    var d = new Date();
    if (d.getTime() > begin + interval) {
      break
    }
    await sleep(100);
  }
  document.body.style.cursor = 'default';
}

updateDisplay();

