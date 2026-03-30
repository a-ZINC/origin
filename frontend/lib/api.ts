const BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// ── types ─────────────────────────────────────────────────────────────────────

export type Post = {
  ID:          string;
  Title:       string;
  Slug:        string;
  Excerpt:     string;
  Body:        string;
  HTMLBody:    string;
  CoverImage:  string;
  Tags:        string[];
  Published:   boolean;
  Visibility:  string;
  ReadingTime: number;
  SeriesID:    string | null;
  SeriesPos:   number | null;
  SeriesName:  string;
  SeriesSlug:  string;
  CreatedAt:   string;
  UpdatedAt:   string;
  PublishedAt: string | null;
  PrevPost:    { Title: string; Slug: string } | null;
  NextPost:    { Title: string; Slug: string } | null;
};

export type Series = {
  ID:          string;
  Name:        string;
  Slug:        string;
  Description: string;
  PostCount:   number;
  CreatedAt:   string;
};

export type Project = {
  ID:          string;
  Name:        string;
  Description: string;
  URL:         string;
  RepoURL:     string;
  Language:    string;
  Stars:       number;
  Featured:    boolean;
  SortOrder:   number;
};

// ── fetchers ──────────────────────────────────────────────────────────────────

async function get<T>(path: string, token?: string): Promise<T> {
  const headers: HeadersInit = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${BASE}${path}`, {
    headers,
    next: { revalidate: 60 },   // ISR — revalidate every 60s
  });

  if (!res.ok) throw new Error(`API ${res.status}: ${path}`);
  return res.json();
}

// ── public API ────────────────────────────────────────────────────────────────

export const getPosts   = (tag?: string) =>
  get<Post[]>(`/api/posts${tag ? `?tag=${tag}` : ""}`);

export const getPost    = (slug: string) =>
  get<Post>(`/api/posts/${slug}`);

export const getSeries  = () =>
  get<Series[]>("/api/series");

export const getSeriesDetail = (slug: string) =>
  get<{ series: Series; posts: Post[] }>(`/api/series/${slug}`);

export const getProjects = () =>
  get<Project[]>("/api/projects");

// ── admin API ─────────────────────────────────────────────────────────────────

export const adminGetPosts = (token: string) =>
  get<Post[]>("/api/admin/posts", token);

export async function adminCreatePost(token: string, data: Partial<Post>) {
  const res = await fetch(`${BASE}/api/admin/posts`, {
    method:  "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
    body:    JSON.stringify(data),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json() as Promise<Post>;
}

export async function adminUpdatePost(token: string, id: string, data: Partial<Post>) {
  const res = await fetch(`${BASE}/api/admin/posts/${id}`, {
    method:  "PUT",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
    body:    JSON.stringify(data),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json() as Promise<Post>;
}

export async function adminDeletePost(token: string, id: string) {
  const res = await fetch(`${BASE}/api/admin/posts/${id}`, {
    method:  "DELETE",
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function adminPreview(token: string, markdown: string) {
  const res = await fetch(`${BASE}/api/posts/preview`, {
    method:  "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
    body:    JSON.stringify({ markdown }),
  });
  if (!res.ok) throw new Error(await res.text());
  const data = await res.json() as { html: string };
  return data.html;
}

export async function trackView(postId: string) {
  await fetch(`${BASE}/api/track`, {
    method:  "POST",
    headers: { "Content-Type": "application/json" },
    body:    JSON.stringify({ postId }),
    keepalive: true,
  }).catch(() => {});   
}