import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { Coffee, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import useUserProfileStore from "@/store/userProfileStore";
import { Input } from "../ui/input";

interface NavbarProps {
  isAuthenticated?: boolean;
  onSignIn?: () => Promise<void>;
}

const Navbar: React.FC<NavbarProps> = ({ isAuthenticated, onSignIn }) => {
  const navigate = useNavigate();
  const { signOut } = useAuth();
  const { profile } = useUserProfileStore();
  const [placeholder, setPlaceholder] = useState("Search");

  useEffect(() => {
    const updatePlaceholder = () => {
      setPlaceholder(
        window.innerWidth >= 640 ? "Search users or coffee shops..." : "Search"
      );
    };

    // Initial call
    updatePlaceholder();

    // Add event listener
    window.addEventListener("resize", updatePlaceholder);

    // Cleanup
    return () => window.removeEventListener("resize", updatePlaceholder);
  }, []);

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

  const handleSearch = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // Implement search functionality here
    console.log("Search submitted");
  };

  return (
    <nav className="fixed top-0 w-full bg-[#FAFAFA] border-b border-[#E5E7EB] z-50 shadow-sm">
      <div className="max-w-6xl mx-auto px-4 py-3 flex justify-between items-center">
        <div
          className="flex items-center gap-3 cursor-pointer hover:opacity-80 transition-opacity"
          onClick={handleLogoClick}
        >
          <Coffee className="h-7 w-7 text-[#4A3726]" />
          <h1 className="text-2xl font-bold text-[#4A3726] font-serif hidden sm:block">
            Coffeehaus
          </h1>
        </div>

        {isAuthenticated && (
          <form onSubmit={handleSearch} className="flex-grow max-w-md mx-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-[#967259]" />
              <Input
                type="search"
                placeholder={placeholder}
                className="w-full pl-10 pr-4 py-2 border-[#D4B08C] focus:ring-[#967259] text-[#4A3726] placeholder-[#967259] bg-[#F9F6F4]"
              />
            </div>
          </form>
        )}

        <div className="flex items-center gap-3">
          {isAuthenticated ? (
            <>
              {profile?.photo_url && (
                <img
                  src={profile.photo_url}
                  alt="Profile"
                  className="w-8 h-8 rounded-full cursor-pointer hover:opacity-80 transition-opacity border-2 border-[#D4B08C]"
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
                className="border-[#4A3726] text-[#4A3726] hover:bg-[#4A3726] hover:text-[#F9F6F4] transition-colors duration-300 font-semibold px-4 py-2 rounded-full text-sm"
                onClick={handleSignOut}
              >
                Sign Out
              </Button>
            </>
          ) : (
            <Button
              variant="outline"
              className="border-[#4A3726] text-[#4A3726] hover:bg-[#4A3726] hover:text-[#F9F6F4] transition-colors duration-300 font-semibold px-4 py-2 rounded-full text-sm"
              onClick={onSignIn}
            >
              Sign In
            </Button>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
