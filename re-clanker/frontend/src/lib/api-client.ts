/**
 * Type-safe API client for the ClankerLoop backend (Go).
 * No authentication required - all endpoints are public.
 */

const getBackendUrl = () => {
  const url = process.env.NEXT_PUBLIC_API_URL;
  if (!url) {
    throw new Error("NEXT_PUBLIC_API_URL environment variable is not set");
  }
  return url;
};

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: { message: string };
}

/**
 * Makes a GET request to the backend API.
 */
export async function apiGet<T>(
  path: string,
): Promise<T> {
  const res = await fetch(`${getBackendUrl()}${path}`, {
    headers: { "Content-Type": "application/json" },
  });
  
  if (!res.ok) {
    const json: ApiResponse<T> = await res.json();
    throw new Error(json.error?.message || "Backend request failed");
  }
  
  const json: ApiResponse<T> = await res.json();
  if (!json.success) {
    throw new Error(json.error?.message || "Backend request failed");
  }
  return json.data as T;
}

/**
 * Makes a POST request to the backend API.
 */
export async function apiPost<T>(
  path: string,
  body?: object,
): Promise<T> {
  const res = await fetch(`${getBackendUrl()}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: body ? JSON.stringify(body) : undefined,
  });
  
  if (!res.ok) {
    const json: ApiResponse<T> = await res.json();
    throw new Error(json.error?.message || "Backend request failed");
  }
  
  const json: ApiResponse<T> = await res.json();
  if (!json.success) {
    throw new Error(json.error?.message || "Backend request failed");
  }
  return json.data as T;
}
