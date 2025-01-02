/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useState } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useAuth } from "@/contexts/AuthContext";
import { mockFeedPosts } from "@/utils/mockData";
import PostCard from "@/components/landing/PostCard";
import WelcomeBanner from "@/components/landing/WelcomeBanner";
import BottomCTA from "@/components/landing/BottomCTA";
import Navbar from "@/components/layout/Navbar";
import { Coffee, MapPin, Users } from "lucide-react";

const LandingPage: React.FC = () => {
  const { signInWithGoogle } = useAuth();
  const [showAuthPrompt, setShowAuthPrompt] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleInteraction = () => {
    setShowAuthPrompt(true);
    setTimeout(() => {
      setShowAuthPrompt(false);
    }, 3000);
  };

  const handleSignIn = async () => {
    try {
      await signInWithGoogle();
    } catch (error) {
      setError("Failed to sign in with Google. Please try again.");
      setTimeout(() => setError(null), 3000);
    }
  };

  return (
    <div className="min-h-screen bg-[#F9F6F4]">
      <Navbar isAuthenticated={false} onSignIn={handleSignIn} />

      {/* Auth Prompt Alert */}
      {showAuthPrompt && (
        <div className="fixed top-16 left-0 right-0 mx-auto w-full max-w-md z-50 px-4">
          <Alert className="bg-[#4A3726] text-[#F9F6F4] border-none shadow-md">
            <AlertDescription>
              Sign up or log in to like, comment, and share your own coffee
              experiences!
            </AlertDescription>
          </Alert>
        </div>
      )}

      {/* Error Alert */}
      {error && (
        <div className="fixed top-16 left-0 right-0 mx-auto w-full max-w-md z-50 px-4">
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        </div>
      )}

      {/* Main Content */}
      <main className="pt-14">
        <div className="max-w-5xl mx-auto px-4 py-8 grid grid-cols-1 md:grid-cols-12 gap-6 lg:gap-8">
          {/* Welcome Banner - Left Side */}
          <div className="md:col-span-4 md:sticky md:top-24 h-fit space-y-6">
            <WelcomeBanner onSignUp={handleSignIn} />

            {/* Features Section */}
            <div className="bg-white p-6 rounded-lg shadow-sm">
              <h3 className="text-xl font-serif font-bold text-[#4A3726] mb-4">
                Why Join Coffeehaus?
              </h3>
              <ul className="space-y-4">
                <li className="flex items-center">
                  <Coffee className="w-5 h-5 text-[#967259] mr-2" />
                  <span className="text-sm text-[#634832]">
                    Discover new coffee shops and brews
                  </span>
                </li>
                <li className="flex items-center">
                  <Users className="w-5 h-5 text-[#967259] mr-2" />
                  <span className="text-sm text-[#634832]">
                    Connect with fellow coffee enthusiasts
                  </span>
                </li>
                <li className="flex items-center">
                  <MapPin className="w-5 h-5 text-[#967259] mr-2" />
                  <span className="text-sm text-[#634832]">
                    Find the best coffee spots near you
                  </span>
                </li>
              </ul>
            </div>
          </div>

          {/* Feed - Right Side */}
          <div className="md:col-span-8">
            <h2 className="text-2xl font-serif font-bold text-[#4A3726] mb-4">
              Featured Posts
            </h2>
            <div className="bg-white shadow-sm rounded-lg overflow-hidden">
              {mockFeedPosts.map((post) => (
                <PostCard
                  key={post.id}
                  post={post}
                  onInteraction={handleInteraction}
                />
              ))}
              <BottomCTA onSignUp={handleSignIn} />
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default LandingPage;
