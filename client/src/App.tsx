import { AuthProvider } from "@/contexts/AuthContext";
import { BrowserRouter, Routes, Route, Navigate } from "react-router";
import { useAuth } from "@/contexts/AuthContext";
import LandingPage from "@/pages/LandingPage";
import ProfileSetupPage from "@/pages/ProfileSetupPage";
import FeedPage from "@/pages/FeedPage";

// Protected route wrapper
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { user, loading } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!user) {
    return <Navigate to="/" />;
  }

  return <>{children}</>;
};

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route
            path="/profile-setup"
            element={
              <ProtectedRoute>
                <ProfileSetupPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/feed"
            element={
              <ProtectedRoute>
                <FeedPage />
              </ProtectedRoute>
            }
          />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
