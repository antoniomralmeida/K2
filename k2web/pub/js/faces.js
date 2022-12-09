const faceWrapper = document.getElementById("face");
const jsonElement = document.getElementById("json");

let face;
const params = new URL(location.href).searchParams;
const avatar = params.get('avatar');

const wait = (ms) =>{
  var dt = new Date();
  var end = dt.getTime() + ms; 
  while (dt.getTime() < end) {
    dt = new Date();
  }
}

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

const randomizeFace = (oldFace, newFace) => {
  Array.from(document.getElementsByClassName("random-attribute")).forEach(
    (elem) => {
      if (!elem.checked) {
        const parts = elem.id.split("-").slice(1);
        if (parts.length === 1) newFace[parts[0]] = oldFace[parts[0]];
        else if (!isNaN(parseInt(parts[1]))) {
          const idx = parseInt(parts[1]);
          newFace[parts[0]][idx] = oldFace[parts[0]][idx];
        } else if (parts.length === 2)
          newFace[parts[0]][parts[1]] = oldFace[parts[0]][parts[1]];
      }
    }
  );
  return newFace;
};

const updateDisplay = () => {  
  var dt = new Date()
  console.log(dt.getTime())
  window.setTimeout(faces.display(faceWrapper, face), 500);
  //wait(500);
};

const SS = () => {
  for (i=0;i<15;i++) {
    speaking();
  }
}

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


const Speak = (text) => {
  interval = text.length * 500
  var d = new Date();
  var begin = d.getTime()
  while (true) {
    speaking();
    var d = new Date();
    if (d.getTime() > begin + interval) {
      break
    }
  }
}

updateDisplay();


