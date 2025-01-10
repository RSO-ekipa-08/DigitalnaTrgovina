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
        <div class="bg-white rounded-lg shadow-lg overflow-hidden hover:shadow-xl transition-shadow duration-300 cursor-pointer" onclick="location.href='/app/${app.id}'">
          <div class="aspect-w-1 aspect-h-1">
            <img src="${app.iconUrl || "/public/images/default-app-icon.png"}"
                 alt="${app.name}"
                 class="w-full h-full object-cover"
                 onerror="this.src='/public/images/default-app-icon.png'">
          </div>
          <div class="p-4">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">${app.name || "Unnamed App"}</h3>
            <div class="space-y-2">
              <div class="text-sm text-gray-600">${app.category || "Uncategorized"}</div>
              <div class="flex justify-between items-center">
                <span class="text-sm font-medium ${app.price > 0 ? 'text-indigo-600' : 'text-green-600'}">${app.price > 0 ? app.price + " €" : "Brezplačno"}</span>
                <span class="flex items-center text-sm text-yellow-500">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                  </svg>
                  ${(app.rating || 0).toFixed(1)}
                </span>
              </div>
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
        <div class="bg-white rounded-lg shadow-lg overflow-hidden hover:shadow-xl transition-shadow duration-300 cursor-pointer" onclick="location.href='/app/${app.id}'">
          <div class="aspect-w-1 aspect-h-1">
            <img src="${app.iconUrl || "/public/images/default-app-icon.png"}"
                 alt="${app.name}"
                 class="w-full h-full object-cover"
                 onerror="this.src='/public/images/default-app-icon.png'">
          </div>
          <div class="p-4">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">${app.name || "Unnamed App"}</h3>
            <div class="space-y-2">
              <div class="text-sm text-gray-600">${app.category || "Uncategorized"}</div>
              <div class="flex justify-between items-center">
                <span class="text-sm font-medium ${app.price > 0 ? 'text-indigo-600' : 'text-green-600'}">${app.price > 0 ? app.price + " €" : "Brezplačno"}</span>
                <span class="flex items-center text-sm text-yellow-500">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                  </svg>
                  ${(app.rating || 0).toFixed(1)}
                </span>
              </div>
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
