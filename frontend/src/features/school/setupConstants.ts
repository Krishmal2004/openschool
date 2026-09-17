export type ALStreamKey = "science_physical" | "science_bio" | "commerce" | "arts" | "technology";

export interface ALStreamDef {
  key: ALStreamKey;
  label: string;
  streamName: string;
  groupName: string | null;
  defaultCode: string;
}

export type AlStreamsState = Record<ALStreamKey, { enabled: boolean; code: string; sections: number }>;

export interface SchoolFormState {
  name: string;
  address: string;
  phone: string;
  email: string;
  logo_url: string;
  school_type: "boys" | "girls" | "mixed";
  grade_from: number | "";
  grade_to: number | "";
}

export const AL_STREAM_DEFS: ALStreamDef[] = [
  { key: "science_physical", label: "Physical Science", streamName: "Science", groupName: "Physical Science", defaultCode: "M" },
  { key: "science_bio", label: "Bio Science", streamName: "Science", groupName: "Bio Science", defaultCode: "B" },
  { key: "commerce", label: "Commerce", streamName: "Commerce", groupName: null, defaultCode: "C" },
  { key: "arts", label: "Arts", streamName: "Arts", groupName: null, defaultCode: "A" },
  { key: "technology", label: "Technology", streamName: "Technology", groupName: null, defaultCode: "T" },
];

export const AL_GRADE_NUMBERS = new Set([12, 13]);

export const GRADE_MIN = 1;
export const GRADE_MAX = 13;
// Mediums come before Classes so generated sections can be tagged with a medium; Rooms is independent and goes last.
export const STEPS = ["School", "Houses", "Grades", "Mediums", "Classes", "Rooms", "Done"] as const;
export const HOUSE_COLOR_PALETTE = ["#0f62fe", "#da1e28", "#24a148", "#f1c21b", "#8a3ffc", "#ff832b"];
export const SUGGESTED_MEDIUMS = ["Sinhala", "Tamil", "English"];

// Special rooms created as "eca" during setup; an admin can retag them as subject Labs later.
export const SUGGESTED_ROOMS = ["Library", "Music Room", "IT Room", "Science Lab", "Auditorium"];
