import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface WelcomeBannerProps {
  onSignUp: () => void;
}

const WelcomeBanner: React.FC<WelcomeBannerProps> = ({ onSignUp }) => {
  return (
    <Card className="mb-6 bg-gradient-to-r from-amber-50 to-orange-50 border-none">
      <CardContent className="pt-6">
        <h2 className="text-2xl font-bold text-amber-900 mb-2">
          Welcome to Coffeehaus
        </h2>
        <p className="text-amber-800 mb-4">
          Join our community of coffee enthusiasts to discover and share the
          best coffee experiences in your area.
        </p>
        <Button
          className="w-full bg-amber-800 text-white hover:bg-amber-700"
          onClick={onSignUp}
        >
          Join Now
        </Button>
      </CardContent>
    </Card>
  );
};

export default WelcomeBanner;
