const appId = window.location.pathname.split("/").pop();

if (!appId) {
  console.error("No app ID found in URL");
  document.body.innerHTML = '<p class="error-message">Invalid app ID</p>';
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
      throw new Error("No app data received");
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
      elements.icon.src =
        app.screenshots?.[0] || "/public/images/default-app-icon.png";
    if (elements.name) elements.name.textContent = app.name || "Unknown App";
    if (elements.developer)
      elements.developer.textContent = app.developerId || "Unknown Developer";
    if (elements.category)
      elements.category.textContent = app.category || "Uncategorized";
    if (elements.price)
      elements.price.textContent =
        app.price > 0 ? `${app.price} €` : "Brezplačno";
    if (elements.description)
      elements.description.textContent =
        app.description || "No description available";

    // Handle screenshots
    if (elements.screenshots) {
      const screenshots = app.screenshots || [];
      elements.screenshots.innerHTML =
        screenshots.length > 0
          ? screenshots
              .map((url) => `<img src="${url}" alt="Screenshot">`)
              .join("")
          : "<p>No screenshots available</p>";
    }
  } catch (error) {
    console.error("Error loading app details:", error);
    const appContent = document.getElementById("app-content");
    if (appContent) {
      appContent.innerHTML =
        '<p class="error-message">Failed to load app details. Please try again later.</p>';
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
          .map((review) => {
            const date = new Date(review.createdAt.split(" ")[0]);
            return `
                      <div class="review">
                          <div class="review-header">
                              <div class="star-rating">
                                  ${"⭐".repeat(review.score)}
                              </div>
                              <div>${date.toLocaleDateString("sl-SI")}</div>
                          </div>
                          <div class="review-comment">${review.comment}</div>
                          <div class="review-author">
                              - ${review.userId}
                          </div>
                      </div>
                  `;
          })
          .join("");
      } else {
        reviewsList.innerHTML =
          "<p>Še ni ocen. Bodite prvi, ki bo ocenil aplikacijo!</p>";
      }
    }
  } catch (error) {
    console.error("Error loading reviews:", error);
    const reviewsList = document.getElementById("reviews-list");
    if (reviewsList) {
      reviewsList.innerHTML =
        '<p class="error-message">Napaka pri nalaganju ocen. Prosimo, poskusite kasneje.</p>';
    }
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
      '<p class="error-message">Application configuration error. Please contact support.</p>';
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
    alert("Prosim, prijavite se za oddajo ocene");
    return;
  }

  if (!comment) {
    alert("Prosim, vnesite komentar");
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
        // Remove variables since we're using inline values
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
      throw new Error("Napaka pri oddaji ocene");
    }

    // Clear the form
    document.getElementById("rating-input").value = "5";
    document.getElementById("review-comment").value = "";

    // Reload reviews to show the new one
    await loadReviews();

    alert("Ocena je bila uspešno oddana");
  } catch (error) {
    console.error("Error submitting review:", error);
    alert("Napaka pri oddaji ocene: " + error.message);
  }
}
