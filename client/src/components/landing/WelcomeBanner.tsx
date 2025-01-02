import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface WelcomeBannerProps {
  onSignUp: () => void;
}

const WelcomeBanner: React.FC<WelcomeBannerProps> = ({ onSignUp }) => {
  return (
    <Card className="mb-6 bg-gradient-to-r from-[#E8D7C9] to-[#F3E9E2] border-none shadow-md overflow-hidden">
      <CardContent className="pt-8 pb-8 px-6">
        <h2 className="text-3xl font-bold text-[#4A3726] mb-4 font-serif">
          Welcome to Coffeehaus
        </h2>
        <p className="text-[#634832] mb-6 text-lg">
          Join our community of coffee enthusiasts to discover and share the
          best coffee experiences in your area.
        </p>
        <Button
          className="w-full bg-[#4A3726] text-[#F9F6F4] hover:bg-[#634832] transition-colors duration-300 text-lg py-3 rounded-full font-semibold"
          onClick={onSignUp}
        >
          Join Now
        </Button>
      </CardContent>
    </Card>
  );
};

export default WelcomeBanner;
