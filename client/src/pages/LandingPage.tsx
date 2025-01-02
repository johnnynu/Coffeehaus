/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useState } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useAuth } from "@/contexts/AuthContext";
import { mockFeedPosts } from "@/utils/mockData";
import PostCard from "@/components/landing/PostCard";
import WelcomeBanner from "@/components/landing/WelcomeBanner";
import BottomCTA from "@/components/landing/BottomCTA";
import Navbar from "@/components/layout/Navbar";

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
    <div className="min-h-screen bg-gray-50">
      <Navbar isAuthenticated={false} onSignIn={handleSignIn} />

      {/* Auth Prompt Alert */}
      {showAuthPrompt && (
        <div className="fixed top-16 left-0 right-0 mx-auto w-full max-w-md z-50 px-4">
          <Alert className="bg-amber-800 text-white border-none">
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
        <div className="max-w-5xl mx-auto px-4 py-6 grid grid-cols-1 md:grid-cols-12 gap-8">
          {/* Welcome Banner - Left Side */}
          <div className="md:col-span-4 md:sticky md:top-20 h-fit">
            <WelcomeBanner onSignUp={handleSignIn} />
          </div>

          {/* Feed - Right Side */}
          <div className="md:col-span-8">
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
      </main>
    </div>
  );
};

export default LandingPage;
