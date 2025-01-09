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
    console.log("Starting loadApps request");
    console.log("APP_SERVICE_URL:", window.APP_SERVICE_URL);

    const requestBody = {
      query: "",
      category: "",
      minPrice: 0,
      maxPrice: 100,
      minAndroidVersion: "",
      pagination: {
        page: 1,
        pageSize: 20,
      },
      sort: [
        {
          field: "name",
          ascending: true,
        },
      ],
    };

    console.log("Request body:", JSON.stringify(requestBody, null, 2));

    const response = await fetch(
      `${window.APP_SERVICE_URL}/app.v1.ApplicationService/SearchApplications`,
      {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
          "Connect-Protocol-Version": "1",
        },
        body: JSON.stringify(requestBody),
        mode: "cors",
        signal: AbortSignal.timeout(5000), // 5 second timeout
      },
    );

    console.log("Response received:", {
      status: response.status,
      statusText: response.statusText,
      headers: Object.fromEntries(response.headers.entries()),
    });

    if (!response.ok) {
      throw new Error(
        `HTTP error! status: ${response.status} statusText: ${response.statusText}`,
      );
    }

    const data = await response.json();
    console.log("Response data:", data);

    if (!data || !data.applications) {
      document.getElementById("apps-grid").innerHTML =
        '<p class="error-message">No applications available.</p>';
      return;
    }

    const appsGrid = document.getElementById("apps-grid");
    appsGrid.innerHTML = data.applications
      .map(
        (app) => `
        <div class="app-card" onclick="location.href='/app/${app.id}'">
          <img src="${app.screenshots?.[0] || "/public/images/default-app-icon.png"}"
               alt="${app.name}">
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

// Add the search and filter functionality
async function applyFilters() {
  const category = document.getElementById("category-filter").value;
  const minPrice = document.getElementById("min-price").value;
  const maxPrice = document.getElementById("max-price").value;

  try {
    const requestBody = {
      query: "",
      category: category || "",
      minPrice: minPrice ? parseFloat(minPrice) : 0,
      maxPrice: maxPrice ? parseFloat(maxPrice) : 100,
      minAndroidVersion: "",
      pagination: {
        page: 1,
        pageSize: 20,
      },
      sort: [
        {
          field: "name",
          ascending: true,
        },
      ],
    };

    const response = await fetch(
      `${window.APP_SERVICE_URL}/app.v1.ApplicationService/SearchApplications`,
      {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
          "Connect-Protocol-Version": "1",
        },
        body: JSON.stringify(requestBody),
        mode: "cors",
        signal: AbortSignal.timeout(5000),
      },
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    const appsGrid = document.getElementById("apps-grid");

    if (!data || !data.applications || data.applications.length === 0) {
      appsGrid.innerHTML =
        '<p class="error-message">No applications found matching your criteria.</p>';
      return;
    }

    appsGrid.innerHTML = data.applications
      .map(
        (app) => `
        <div class="app-card" onclick="location.href='/app/${app.id}'">
          <img src="${app.screenshots?.[0] || "https://developer.android.com/static/images/logos/android.svg"}"
               alt="${app.name}"
               onerror="this.src='https://developer.android.com/static/images/logos/android.svg'">
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

// Load categories
async function loadCategories() {
  try {
    const response = await fetch(
      `${window.APP_SERVICE_URL}/app.v1.ApplicationService/ListCategories`,
      {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          pagination: {
            page: 1,
            page_size: 40,
          },
        }),
      },
    );

    // Log the full response for debugging
    console.log("Response status:", response.status);
    console.log(
      "Response headers:",
      Object.fromEntries(response.headers.entries()),
    );
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();

    if (!data || !data.categories) {
      console.warn("No categories available");
      return;
    }

    // Populate category dropdown
    const categorySelect = document.getElementById("category-filter");
    const options = data.categories.map(
      (category) => `<option value="${category.id}">${category.name}</option>`,
    );

    // Keep the default "All categories" option and add new ones
    categorySelect.innerHTML = `
      <option value="">Vse kategorije</option>
      ${options.join("")}
    `;

    console.log("Categories loaded successfully");
  } catch (error) {
    console.error("Error loading categories:", error);
  }
}

// Initialize
document.addEventListener("DOMContentLoaded", async () => {
  if (window.APP_SERVICE_URL && window.REVIEWS_SERVICE_URL) {
    updateAuthUI();
    await loadCategories(); // Load categories first
    await loadApps(); // Then load apps
  } else {
    console.error("Service URLs not configured");
    document.body.innerHTML =
      '<p class="error-message">Application configuration error. Please contact support.</p>';
  }
});
