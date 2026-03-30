"use client";
import { useEffect } from "react";
import { trackView }  from "@/lib/api";

export default function ViewTracker({ postId }: { postId: string }) {
  useEffect(() => {
    trackView(postId);
  }, [postId]);

  return null; // renders nothing — side effect only
}
