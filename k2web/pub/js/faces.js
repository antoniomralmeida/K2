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
  console.log(face);
  faces.display(faceWrapper, face);
  history.replaceState(
    undefined,
    undefined,
    `?avatar=${btoa(JSON.stringify(face))}`
  );
  //jsonElement.value = JSON.stringify(face);
};



const isValue = (obj) =>
  typeof obj === "boolean" ||
  typeof obj === "number" ||
  typeof obj === "string";


const getValue = (oldValue, event) => {
  if (typeof oldValue === "number") {
    return parseFloat(event.target.value);
  }
  if (typeof oldValue === "boolean") {
    return event.target.checked;
  }
  return event.target.value;
};
const listenForChanges = () => {
  const textInputs = document.querySelectorAll("input.form-control");
  const checkboxInputs = document.querySelectorAll(
    "input.form-check-input"
  );
  const selects = document.querySelectorAll("select.form-control");

  for (const input of [...textInputs, ...checkboxInputs, ...selects]) {
    if (input.id.startsWith("randomize")) continue;
    input.addEventListener("change", (event) => {
      const parts = event.target.id.split("-");

      if (parts.length === 1) {
        face[parts[0]] = getValue(face[parts[0]], event);
      } else if (parts.length === 2) {
        face[parts[0]][parts[1]] = getValue(
          face[parts[0]][parts[1]],
          event
        );
      } else {
        throw new Error(`Invalid ID ${event.target.id}`);
      }

      updateDisplay();
    });
  }


  const checkboxes = [].concat(
    Array.from(document.getElementsByClassName("random-group")),
    Array.from(document.getElementsByClassName("random-attribute"))
  );



  Array.from(document.getElementsByClassName("random-group")).forEach(
    (elem) => {
      elem.addEventListener("click", () => {
        Array.from(
          document.getElementsByClassName(elem.id.replace("-group", ""))
        ).forEach((elem2) => (elem2.checked = elem.checked));
      });
    }
  );

  Array.from(document.getElementsByClassName("random-attribute")).forEach(
    (elem) => {
      elem.addEventListener("click", () => {
        Array.from(elem.className.split(" ")).forEach((className) => {
          if (className.startsWith("randomize-")) {
            if (elem.checked) {
              // See if all are checked, and if so check the group
              const others = Array.from(
                document.getElementsByClassName(className)
              ).filter(
                (other) => !other.classList.contains("random-group")
              );
              if (others.every((other) => other.checked)) {
                document.getElementById(
                  className + "-group"
                ).checked = true;
              }
            } else {
              document.getElementById(className + "-group").checked = false;
            }
          }
        });
      });
    }
  );
};


updateDisplay();
listenForChanges();

