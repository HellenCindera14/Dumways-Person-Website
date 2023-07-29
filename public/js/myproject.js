function emptyFormAlert() {
  let name = document.getElementById("input-name-project").value;
  let startDate = document.getElementById("start-date").value;
  let endDate = document.getElementById("end-date").value;
  let description = document.getElementById("input-description").value;
  let multiInput = document.querySelectorAll(".multi-input:checked");
  let image = document.getElementById("input-image").value;

  if (name == "") {
    return alert("Please input your project name or title");
  } else if (startDate == "") {
    return alert("When did you start this project?");
  } else if (endDate == "") {
    return alert("When did you finish this project?");
  } else if (description == "") {
    return alert("Please describe this project.");
  } else if (multiInput.length === 0) {
    return alert("Choose your technologies.");
  } else if (image == "") {
    return alert("Please attach an image of your project.");
  }
}

let dataProject = [];

function addProject(event) {
  event.preventDefault();
  let name = document.getElementById("input-name-project").value;
  let startDate = document.getElementById("start-date").value;
  let endDate = document.getElementById("end-date").value;
  let description = document.getElementById("input-description").value;
  let image = document.getElementById("input-image").files;

  const nodeJsIcon = '<i class="fa-brands fa-node-js fs-3"></i>';
  const golangIcon = '<i class="fa-brands fa-golang fs-3"></i>';
  const reactJsIcon = '<i class="fa-brands fa-react fs-3"></i>';
  const javascriptIcon = '<i class="fa-brands fa-square-js fs-3"></i>';

  let nodejs = document.getElementById("node-js").checked ? nodeJsIcon : "";
  let golang = document.getElementById("golang").checked ? golangIcon : "";
  let reactjs = document.getElementById("react-js").checked ? reactJsIcon : "";
  let javascript = document.getElementById("javascript").checked
    ? javascriptIcon
    : "";

  // membuat url gambar dan menampilkan gambar yg dipilih pertama
  image = URL.createObjectURL(image[0]);
  console.log(image);

  let mulai = new Date(startDate);
  let akhir = new Date(endDate);

  if (mulai > akhir) {
    return alert("Please input your dates correctly.");
  }

  let selisih = akhir.getTime() - mulai.getTime();
  let days = selisih / (1000 * 60 * 60 * 24);
  let weeks = Math.floor(days / 7);
  let months = Math.floor(weeks / 4);
  let years = Math.floor(months / 12);
  let durasi = "";

  if (days < 7) {
    durasi = days + " Hari";
  } else if (days >= 7 && weeks < 4) {
    durasi = weeks + " Minggu";
  } else if (weeks >= 4 && months <= 12) {
    durasi = months + " Bulan";
  } else {
    durasi = years + " Tahun";
  }

  let data = {
    name,
    days,
    weeks,
    months,a
    years,
    selisih,
    durasi,
    description,
    nodejs,
    golang,
    reactjs,
    javascript,
    image,
  };

  dataProject.push(data);
  console.log(dataProject);

  renderDataProject();
}

function renderDataProject() {
  document.getElementById("contents").innerHTML = "";

  for (let index = 0; index < dataProject.length; index++) {
    document.getElementById("contents").innerHTML += `

    <div class="card card-project">
          <img src="${dataProject[index].image}" alt="" />
          <div class="card-title">
            <p class="card-title fs-6 fw-bold mb-0">
            <a href="blog-detail" target="_blank"
            >${dataProject[index].name}</a>
            </p>
            <div class="card-duration">durasi: ${dataProject[index].durasi}</div>
            <p class="card-text my-2">
            ${dataProject[index].description}
            </p>
            <div class="d-flex gap-4 my-3" style="width: 100%">
            ${dataProject[index].nodejs}
            ${dataProject[index].golang}
            ${dataProject[index].reactjs}
            ${dataProject[index].javascript}
            </div>
            <div class="btn-card-index d-flex gap-2" style="width: 100%">
              <button class="btn btn-dark" style="width: 100%">edit</button>
              <button class="btn btn-dark" style="width: 100%">delete</button>
            </div>
          </div>
    </div>

    
    `;
  }
}

{
  /* <div class="card-project">
          <div class="image-project">
            <img src="${dataProject[index].image}" alt="" />
          </div>

          <div class="duration">
            <a href="blog.html" target="_blank"
              ><h4>${dataProject[index].name}</h4></a
            >
            <p>duration : ${dataProject[index].durasi}</p>
          </div>

          <div class="description">
            <p>
              ${dataProject[index].description}
            </p>
          </div>

          <div class="icons">
            ${dataProject[index].nodejs}
            ${dataProject[index].golang}
            ${dataProject[index].reactjs}
            ${dataProject[index].javascript}
          </div>

          <div class="button-project">
            <button>edit</button>
            <button>delete</button>
          </div>
        </div> */
}

{
  /* <div class="col">
    <div class="card">
      <img src="${dataProject[index].image}" class="card-img-top" alt="" />
      <div class="card-body">
        <h5 class="card-title">${dataProject[index].name}</h5>
        <div class="mb-2 fw-light">durasi: ${dataProject[index].durasi}</div>
        <p class="card-text">
        ${dataProject[index].description}
        </p>
                <div class="d-flex gap-4 mb-3" style="width: 100%">
                      ${dataProject[index].nodejs}
                      ${dataProject[index].golang}
                      ${dataProject[index].reactjs}
                      ${dataProject[index].javascript}
                </div>
        <div class="d-flex gap-2" style="width: 100%">
          <button class="btn btn-dark" style="width: 100%">edit</button>
          <button class="btn btn-dark" style="width: 100%">
            delete
          </button>
        </div>
      </div>
    </div>
  </div> */
}
