import { UserGeneratedVideos } from "./components/UserGeneratedVideos";

interface YourVideosProps {
  onVideoSelect: (videoUrl: string) => void;
}

export const YourVideos = ({ onVideoSelect }: YourVideosProps) => {
  return <UserGeneratedVideos onVideoSelect={onVideoSelect} />;
};
