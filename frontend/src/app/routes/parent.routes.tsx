import { Route } from "react-router";
import { lazy as page } from "react";
import ProtectedRoute from "@/shared/auth/ProtectedRoute";
import ParentLayout from "@/layouts/ParentLayout";

const NotFound = page(() => import("@/app/pages/NotFound"));
const ParentDashboard = page(() => import("@/features/parent/pages/ParentDashboard"));
const ChildDetail = page(() => import("@/features/parent/pages/ChildDetail"));
const NotificationCenter = page(() => import("@/features/notifications/pages/NotificationCenter"));

export function parentRoutes() {
  return (
    <Route element={<ProtectedRoute><ParentLayout /></ProtectedRoute>}>
      <Route index element={<ParentDashboard />} />
      <Route path="/p/children/:id" element={<ChildDetail />} />
      <Route path="/notification-center" element={<NotificationCenter />} />
      <Route path="*" element={<NotFound />} />
    </Route>
  );
}
