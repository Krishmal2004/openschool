import { useQuery } from "@tanstack/react-query";
import { studentSelfApi } from "@/features/students/api/studentSelf";
import { studentKeys } from "@/features/students/keys";

export const useMyStudentProfile = () => useQuery({ queryKey: studentKeys.me.profile(), queryFn: studentSelfApi.profile });

export const useMyAttendance = () => useQuery({ queryKey: studentKeys.me.attendance(), queryFn: studentSelfApi.attendance });

export const useMyMarks = (termId: string) =>
  useQuery({ queryKey: studentKeys.me.marks(termId), queryFn: () => studentSelfApi.marks(termId), enabled: !!termId });
