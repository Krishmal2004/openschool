import api from "@/shared/api/client";

export interface SearchResultItem {
  id: string;
  name: string;
  subtitle: string;
}

// Matches backend GlobalSearchResponse: top 5 per entity for the header jump-to-record search.
export interface GlobalSearchResponse {
  students: SearchResultItem[];
  teachers: SearchResultItem[];
  guardians: SearchResultItem[];
  non_academic_staff: SearchResultItem[];
}

export const searchApi = {
  global: (q: string) =>
    api.get<GlobalSearchResponse>("/admin/search", { params: { q } }).then((r) => r.data),
};
