/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useState } from "react";
import { Coffee } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useAuth } from "@/contexts/AuthContext";
import { mockFeedPosts } from "@/utils/mockData";
import PostCard from "@/components/landing/PostCard";
import WelcomeBanner from "@/components/landing/WelcomeBanner";
import BottomCTA from "@/components/landing/BottomCTA";

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
      {/* Fixed Navigation */}
      <nav className="fixed top-0 w-full bg-white border-b z-50">
        <div className="max-w-5xl mx-auto px-4 py-2 flex justify-between items-center">
          <div className="flex items-center gap-2">
            <Coffee className="h-6 w-6 text-amber-800" />
            <h1 className="text-xl font-bold text-amber-800">Coffeehaus</h1>
          </div>
          <Button
            variant="outline"
            className="border-amber-800 text-amber-800 hover:bg-amber-50"
            onClick={handleSignIn}
          >
            Sign In with Google
          </Button>
        </div>
      </nav>

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
