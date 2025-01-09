// Check authentication status and update UI
function updateAuthUI() {
  const profile = auth.getProfile();
  if (profile) {
    document.getElementById("login-btn").style.display = "none";
    document.getElementById("user-info").style.display = "flex";
    document.getElementById("user-name").textContent = profile.name;
    document.getElementById("user-avatar").src = profile.picture;
  }
}

async function loadApps() {
  try {
    const response = await fetch(
      `${window.APP_SERVICE_URL}/app.v1.ApplicationService/SearchApplications`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Connect-Protocol-Version": "1",
        },
        body: JSON.stringify({
          pagination: {
            page: 1,
            page_size: 20,
          },
          include_moderated_only: true,
        }),
      },
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();

    if (!data || !data.applications) {
      throw new Error("Invalid data format received from server");
    }

    const appsGrid = document.getElementById("apps-grid");
    if (data.applications.length === 0) {
      appsGrid.innerHTML = "<p>No applications found.</p>";
      return;
    }

    appsGrid.innerHTML = data.applications
      .map(
        (app) => `
            <div class="app-card" onclick="location.href='/app/${app.id}'">
                <img src="${app.screenshots?.[0] || "/public/images/default-app-icon.png"}"
                     alt="${app.name}"
                     onerror="this.src='/public/images/default-app-icon.png'">
                <div class="app-card-content">
                    <h3>${app.name || "Unnamed App"}</h3>
                    <div class="app-meta">
                        <div>${app.category || "Uncategorized"}</div>
                        <div>${app.price > 0 ? app.price + " €" : "Brezplačno"}</div>
                        <div>⭐ ${(app.rating || 0).toFixed(1)}</div>
                    </div>
                </div>
            </div>
        `,
      )
      .join("");
  } catch (error) {
    console.error("Error loading apps:", error);
    document.getElementById("apps-grid").innerHTML =
      '<p class="error-message">Failed to load applications. Please try again later.</p>';
  }
}

// Update loadAppDetails in app-details.js
async function loadAppDetails() {
  try {
    const response = await fetch(
      `${window.APP_SERVICE_URL}/app.v1.ApplicationService/GetApplication`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Connect-Protocol-Version": "1",
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

    document.getElementById("app-icon").src =
      app.screenshots?.[0] || "/public/images/default-app-icon.png";
    document.getElementById("app-name").textContent = app.name || "Unknown App";
    document.getElementById("app-developer").textContent =
      app.developerId || "Unknown Developer";
    document.getElementById("app-category").textContent =
      app.category || "Uncategorized";
    document.getElementById("app-price").textContent =
      app.price > 0 ? `${app.price} €` : "Brezplačno";
    document.getElementById("app-description").textContent =
      app.description || "No description available";

    // Handle screenshots
    const screenshots = app.screenshots || [];
    document.getElementById("app-screenshots").innerHTML =
      screenshots.length > 0
        ? screenshots
            .map((url) => `<img src="${url}" alt="Screenshot">`)
            .join("")
        : "<p>No screenshots available</p>";
  } catch (error) {
    console.error("Error loading app details:", error);
    document.getElementById("app-content").innerHTML =
      '<p class="error-message">Failed to load app details. Please try again later.</p>';
  }
}

// Add the search and filter functionality
async function applyFilters() {
  const category = document.getElementById("category-filter").value;
  const minPrice = document.getElementById("min-price").value;
  const maxPrice = document.getElementById("max-price").value;

  try {
    const response = await fetch(
      `${window.APP_SERVICE_URL}/app.v1.ApplicationService/SearchApplications`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Connect-Protocol-Version": "1",
        },
        body: JSON.stringify({
          category: category || undefined,
          minPrice: minPrice ? parseFloat(minPrice) : undefined,
          maxPrice: maxPrice ? parseFloat(maxPrice) : undefined,
          pagination: {
            page: 1,
            page_size: 20,
          },
        }),
      },
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    // Update the UI with filtered results
    let appsGrid = document.getElementById("apps-grid");
    appsGrid.innerHTML = data.applications
      .map(
        (app) => `
            <div class="app-card" onclick="location.href='/app/${app.id}'">
                <img src="${app.screenshots?.[0] || "/public/images/default-app-icon.png"}"
                     alt="${app.name}"
                     onerror="this.src='/public/images/default-app-icon.png'">
                <div class="app-card-content">
                    <h3>${app.name || "Unnamed App"}</h3>
                    <div class="app-meta">
                        <div>${app.category || "Uncategorized"}</div>
                        <div>${app.price > 0 ? app.price + " €" : "Brezplačno"}</div>
                        <div>⭐ ${(app.rating || 0).toFixed(1)}</div>
                    </div>
                </div>
            </div>
        `,
      )
      .join("");
  } catch (error) {
    console.error("Error applying filters:", error);
    document.getElementById("apps-grid").innerHTML =
      '<p class="error-message">Failed to apply filters. Please try again later.</p>';
  }
}

// Initialize
document.addEventListener("DOMContentLoaded", () => {
  updateAuthUI();
  loadApps();
});
