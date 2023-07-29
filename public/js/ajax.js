const promise = new Promise((resolve, reject) => {
  const xhr = new XMLHttpRequest();
  xhr.open("GET", "https://api.npoint.io/b70c5c1e86834bb3f99f", true);
  //   console.log(xhr);
  xhr.onload = () => {
    if (xhr.status === 200) {
      // We parsing it so it is easier to read in console
      // response vs responseText, the differences are, responseText is an older version, when response is more newer, but the output is still the same/similiar.
      resolve(JSON.parse(xhr.response));
    } else {
      reject("Error loading data.");
    }
  };
  xhr.onerror = () => {
    reject("Network error.");
  };
  xhr.send();
});

async function showAll() {
  const cardData = await promise;
  //   console.log(cardData);

  let testimonialHTML = "";
  cardData.forEach((item) => {
    testimonialHTML += `<div class="card-testimonial">
                            <img
                                src="${item.image}"
                                class="profile-testimonial"
                            />
                                <p class="quote">${item.quote}</p>
                                <p class="author">${item.rating} <i class="fa-solid fa-star"></i> from ${item.author}</p>
                            </div>
                          `;
  });

  document.getElementById("container-testimonials").innerHTML = testimonialHTML;
}

// eksekusi awal / default
showAll();

async function stars(rating) {
  const cardData = await promise;

  const testimonialFiltered = cardData.filter((item) => {
    return item.rating === rating;
  });

  //   console.log(testimonialFiltered);

  let testimonialHTML = "";

  if (testimonialFiltered.length === 0) {
    testimonialHTML = "<h1>Data not found!</h1>";
  } else {
    testimonialFiltered.forEach((item) => {
      testimonialHTML += `<div class="card-testimonial">
                                <img
                                    src="${item.image}"
                                    class="profile-testimonial"
                                />
                                <p class="quote">${item.quote}</p>
                                <p class="author">${item.rating} <i class="fa-solid fa-star"></i> from ${item.author}</p>
                            </div>
                          `;
    });
  }

  document.getElementById("container-testimonials").innerHTML = testimonialHTML;
}
