const appId = window.location.pathname.split("/").pop();
let appPrice = 0;

if (!appId) {
  console.error("No app ID found in URL");
  document.body.innerHTML =
    '<div class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 border border-red-200 dark:border-red-800"><p class="text-center text-red-600 dark:text-red-400">Neveljaven ID aplikacije.</p></div>';
}

async function loadAppDetails() {
  if (!appId) return;

  try {
    const response = await fetch(
      `${window.APP_SERVICE_URL}/app.v1.ApplicationService/GetApplication`,
      {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          id: appId,
        }),
      },
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    const app = data.application;

    if (!app) {
      throw new Error("Aplikacija ni bila najdena");
    }

    // Check if elements exist before updating them
    const elements = {
      icon: document.getElementById("app-icon"),
      name: document.getElementById("app-name"),
      developer: document.getElementById("app-developer"),
      category: document.getElementById("app-category"),
      price: document.getElementById("app-price"),
      description: document.getElementById("app-description"),
      screenshots: document.getElementById("app-screenshots"),
    };

    // Update elements if they exist
    if (elements.icon)
      elements.icon.src = app.iconUrl || "/public/images/default-app-icon.png";
    if (elements.name) elements.name.textContent = app.name || "Unknown App";
    if (elements.developer)
      elements.developer.textContent = app.developerId || "Unknown Developer";
    if (elements.category)
      elements.category.textContent = app.category || "Uncategorized";
    if (elements.price) {
      appPrice = app.price; // Store the price
      elements.price.textContent =
        app.price > 0 ? `${app.price} €` : "Brezplačno";
      updateButtonText(app.price);
    }
    if (elements.description)
      elements.description.textContent =
        app.description || "No description available";

    // Handle screenshots
    if (elements.screenshots) {
      const screenshots = app.screenshots || [];
      elements.screenshots.innerHTML =
        screenshots.length > 0
          ? screenshots
              .map(
                (url) =>
                  `<img src="${url}" alt="Zaslonska slika" class="rounded-lg shadow-lg">`,
              )
              .join("")
          : '<div class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 border border-red-200 dark:border-red-800"><p class="text-center text-red-600 dark:text-red-400">Ni razpoložljivih zaslonskih slik.</p></div>';
    }
  } catch (error) {
    console.error("Error loading app details:", error);
    const appContent = document.getElementById("app-content");
    if (appContent) {
      appContent.innerHTML =
        '<div class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 border border-red-200 dark:border-red-800"><p class="text-center text-red-600 dark:text-red-400">Napaka pri nalaganju podrobnosti aplikacije. Prosimo, poskusite kasneje.</p></div>';
    }
  }
}

async function loadReviews() {
  if (!appId) return;

  const query = `
        query GetReviews($appId: String!, $page: Int!, $pageSize: Int!, $includeModeratedOnly: Boolean!, $tenantId: String!) {
            reviews(
                appId: $appId,
                page: $page,
                pageSize: $pageSize,
                includeModeratedOnly: $includeModeratedOnly,
                tenantId: $tenantId
            ) {
                reviews {
                    id
                    userId
                    score
                    comment
                    createdAt
                }
                totalCount
                averageScore
            }
        }
    `;

  try {
    const response = await fetch(`${window.REVIEWS_SERVICE_URL}/graphql`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({
        query,
        variables: {
          appId,
          page: 0,
          pageSize: 10,
          includeModeratedOnly: false,
          tenantId: "default",
        },
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();

    if (result.errors) {
      throw new Error(result.errors[0].message);
    }

    if (!result.data || !result.data.reviews) {
      throw new Error("No review data received");
    }

    const { reviews, totalCount, averageScore } = result.data.reviews;
    console.log(reviews);
    console.log(totalCount);
    console.log(averageScore);

    const averageRatingElement = document.getElementById("average-rating");
    const totalReviewsElement = document.getElementById("total-reviews");
    const reviewsList = document.getElementById("reviews-list");

    if (averageRatingElement) {
      averageRatingElement.textContent = `⭐ ${(averageScore || 0).toFixed(1)}`;
    }

    if (totalReviewsElement) {
      // Proper Slovenian grammar for review count
      if (totalCount === 0) {
        totalReviewsElement.textContent = "Ni ocen";
      } else if (totalCount === 1) {
        totalReviewsElement.textContent = "1 ocena";
      } else if (totalCount === 2) {
        totalReviewsElement.textContent = "2 oceni";
      } else if (totalCount === 3 || totalCount === 4) {
        totalReviewsElement.textContent = `${totalCount} ocene`;
      } else {
        totalReviewsElement.textContent = `${totalCount} ocen`;
      }
    }

    if (reviewsList) {
      if (reviews && reviews.length > 0) {
        console.log(reviews);
        reviewsList.innerHTML = reviews
          .map(
            (review) => `
            <div class="border-t border-gray-200 dark:border-gray-700 pt-4">
              <div class="flex items-center justify-between mb-2">
                <div class="font-medium text-gray-900 dark:text-white">${review.userId}</div>
                <div class="text-sm text-gray-500 dark:text-gray-400">${new Date(
                  review.createdAt.split(" ")[0],
                ).toLocaleString("sl-SI", {
                  year: "numeric",
                  month: "2-digit",
                  day: "2-digit",
                  hour: "2-digit",
                  minute: "2-digit",
                })}</div>
              </div>
              <div class="flex items-center mb-2">
                ${Array(5)
                  .fill()
                  .map(
                    (_, i) =>
                      `<svg class="w-4 h-4 ${
                        i < review.score
                          ? "text-yellow-500 dark:text-yellow-400"
                          : "text-gray-300 dark:text-gray-600"
                      }" fill="currentColor" viewBox="0 0 20 20">
                        <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                      </svg>`,
                  )
                  .join("")}
              </div>
              <p class="text-gray-600 dark:text-gray-400">${review.comment}</p>
            </div>
          `,
          )
          .join("");
      } else {
        reviewsList.innerHTML =
          '<div class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 border border-red-200 dark:border-red-800"><p class="text-center text-red-600 dark:text-red-400">Še ni ocen. Bodite prvi, ki bo ocenil aplikacijo!</p></div>';
      }
    }
  } catch (error) {
    console.error("Error loading reviews:", error);
    const reviewsList = document.getElementById("reviews-list");
    if (reviewsList) {
      reviewsList.innerHTML =
        '<div class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 border border-red-200 dark:border-red-800"><p class="text-center text-red-600 dark:text-red-400">Napaka pri nalaganju ocen. Prosimo, poskusite kasneje.</p></div>';
    }
  }
}

function updateButtonText(price) {
  const button = document.getElementById("download-btn");
  button.textContent = price > 0 ? "Kupi" : "Prenesi";
}

async function handleDownload() {
  const profile = auth.getProfile();
  if (!profile) {
    window.location.href = "/login";
    return;
  }
  console.log(profile);

  const stringToUUID = (str) => {
    // Create a hash from the string
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      const char = str.charCodeAt(i);
      hash = (hash << 5) - hash + char;
      hash = hash & hash; // Convert to 32-bit integer
    }

    // Convert hash to hex string and pad with zeros
    let hex = Math.abs(hash).toString(16).padStart(32, "0");

    // Insert UUID dashes at correct positions
    return (
      hex.substring(0, 8) +
      "-" +
      hex.substring(8, 12) +
      "-4" +
      hex.substring(13, 16) +
      "-a" +
      hex.substring(17, 20) +
      "-" +
      hex.substring(20, 32)
    );
  };

  const userId = stringToUUID(profile.aud);
  console.log("Generated userId:", userId);
  try {
    if (appPrice > 0) {
      // Handle paid app - purchase flow
      const response = await fetch(
        `${window.APP_SERVICE_URL}/app.v1.ApplicationService/BuyApplication`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
            "Connect-Protocol-Version": "1",
          },
          body: JSON.stringify({
            applicationId: appId,
            userId: userId,
          }),
        },
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();

      if (data.checkoutUrl) {
        window.location.href = data.checkoutUrl;
      }
    } else {
      console.log("Sending download request for:", {
        applicationId: appId,
        userId: userId,
      });

      const response = await fetch(
        `${window.APP_SERVICE_URL}/app.v1.ApplicationService/DownloadApplication`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
            "Connect-Protocol-Version": "1",
          },
          body: JSON.stringify({
            applicationId: appId,
            userId: userId,
          }),
        },
      );

      console.log("Response status:", response.status);
      const responseText = await response.text();
      console.log("Response text:", responseText);

      if (!response.ok) {
        throw new Error(
          `HTTP error! status: ${response.status} response: ${responseText}`,
        );
      }

      let data;
      try {
        data = JSON.parse(responseText);
      } catch (e) {
        console.error("Error parsing response:", e);
        throw new Error("Invalid response format");
      }

      if (data.downloadUrl) {
        console.log("Download URL received:", data.downloadUrl);
        const link = document.createElement("a");
        link.href = data.downloadUrl;
        link.setAttribute("download", "");
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      } else {
        throw new Error("No download URL received");
      }
    }
  } catch (error) {
    console.error("Error during download/purchase:", error);
    alert("Napaka pri prenosu/nakupu aplikacije. Prosimo, poskusite kasneje.");
  }
}

// Initialize
document.addEventListener("DOMContentLoaded", () => {
  if (window.APP_SERVICE_URL && window.REVIEWS_SERVICE_URL) {
    updateAuthUI();
    loadAppDetails();
    loadReviews();
  } else {
    console.error("Service URLs not configured");
    document.body.innerHTML =
      '<div class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 border border-red-200 dark:border-red-800"><p class="text-center text-red-600 dark:text-red-400">Napaka pri konfiguraciji aplikacije. Prosimo, kontaktirajte podporo.</p></div>';
  }
});

function updateAuthUI() {
  const profile = auth.getProfile();
  if (profile) {
    document.getElementById("login-btn").style.display = "none";
    document.getElementById("user-info").style.display = "flex";
    document.getElementById("user-name").textContent = profile.name;
    document.getElementById("user-avatar").src = profile.picture;
  } else {
    document.getElementById("login-btn").style.display = "flex";
    document.getElementById("user-info").style.display = "none";
  }
}

async function submitReview() {
  const score = parseInt(document.getElementById("rating-input").value);
  const comment = document.getElementById("review-comment").value.trim();

  const profile = auth.getProfile();
  console.log(profile);
  if (!profile) {
    alert("Prosimo, prijavite se za oddajo ocene.");
    return;
  }

  if (!comment) {
    alert("Prosimo, vnesite komentar.");
    return;
  }

  const mutation = `
        mutation {
            createReview(
                appId: "${appId}",
                userId: "${profile.nickname}",
                score: ${score},
                comment: "${comment}",
                tenantId: "default"
            ) {
                id
                appId
                userId
                score
                comment
                createdAt
                isModerated
                moderationStatus
                tenantId
            }
        }
    `;

  try {
    const response = await fetch(`${window.REVIEWS_SERVICE_URL}/graphql`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        Authorization: `Bearer ${auth.getToken()}`,
      },
      body: JSON.stringify({
        query: mutation,
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();

    if (result.errors) {
      throw new Error(result.errors[0].message);
    }

    if (!result.data?.createReview) {
      throw new Error("Napaka pri oddaji ocene. Prosimo, poskusite ponovno.");
    }

    // Clear the form
    document.getElementById("rating-input").value = "5";
    document.getElementById("review-comment").value = "";

    // Reload reviews to show the new one
    await loadReviews();

    alert("Ocena je bila uspešno oddana. Hvala za vaš komentar!");
  } catch (error) {
    console.error("Error submitting review:", error);
    alert(
      "Napaka pri oddaji ocene. Prosimo, poskusite kasneje: " + error.message,
    );
  }
}
