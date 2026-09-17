import type { PrefectRank } from "@/features/portfolio/api/prefect";

export const PREFECT_RANKS: { value: PrefectRank; label: string }[] = [
  { value: "head", label: "Head Prefects" },
  { value: "deputy_head", label: "Deputy Head Prefects" },
  { value: "senior", label: "Senior Prefects" },
  { value: "junior", label: "Junior Prefects" },
  { value: "house_captain", label: "House Captains" },
  { value: "vice_house_captain", label: "Vice House Captains" },
];
