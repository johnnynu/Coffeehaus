import React from "react";
import { useNavigate } from "react-router";
import { Coffee } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import useUserProfileStore from "@/store/userProfileStore";

interface NavbarProps {
  isAuthenticated?: boolean;
  onSignIn?: () => Promise<void>;
}

const Navbar: React.FC<NavbarProps> = ({ isAuthenticated, onSignIn }) => {
  const navigate = useNavigate();
  const { signOut } = useAuth();
  const { profile } = useUserProfileStore();

  const handleSignOut = async () => {
    try {
      await signOut();
      navigate("/");
    } catch (error) {
      console.error("Error signing out:", error);
    }
  };

  const handleProfileClick = () => {
    if (profile?.username) {
      navigate(`/profile/${profile.username}`);
    }
  };

  const handleLogoClick = () => {
    navigate(isAuthenticated ? "/feed" : "/");
  };

  return (
    <nav className="fixed top-0 w-full bg-white border-b z-50">
      <div className="max-w-5xl mx-auto px-4 py-2 flex justify-between items-center">
        <div
          className="flex items-center gap-2 cursor-pointer hover:opacity-80 transition-opacity"
          onClick={handleLogoClick}
        >
          <Coffee className="h-6 w-6 text-amber-800" />
          <h1 className="text-xl font-bold text-amber-800">Coffeehaus</h1>
        </div>
        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <>
              {profile?.photo_url && (
                <img
                  src={profile.photo_url}
                  alt="Profile"
                  className="w-8 h-8 rounded-full cursor-pointer hover:opacity-80 transition-opacity"
                  crossOrigin="anonymous"
                  referrerPolicy="no-referrer"
                  onClick={handleProfileClick}
                  onError={(e) => {
                    console.error("Error loading profile image:", e);
                    console.log("Attempted image URL:", profile.photo_url);
                  }}
                />
              )}
              <Button
                variant="outline"
                className="border-amber-800 text-amber-800 hover:bg-amber-50"
                onClick={handleSignOut}
              >
                Sign Out
              </Button>
            </>
          ) : (
            <Button
              variant="outline"
              className="border-amber-800 text-amber-800 hover:bg-amber-50"
              onClick={onSignIn}
            >
              Sign In with Google
            </Button>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
