const faceWrapper = document.getElementById("face");
const jsonElement = document.getElementById("json");

let face;
const params = new URL(location.href).searchParams;
const avatar = params.get('avatar');


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
    P1();
  } else {
    P2();
  }
}

const P1 = () => {
  face["mouth"].id = "mouth7";
  updateDisplay();
}

const P2 = () => {
  face["mouth"].id = "smile";
  updateDisplay();
}

const TTS = async(text) => {
  let msg = new SpeechSynthesisUtterance();
  msg.voice = speechSynthesis.getVoices()[1];
  msg.text = text;
  speechSynthesis.speak(msg);
  
}

const Speak = async (text) => {
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


