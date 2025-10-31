import { Title } from "@/components/dashboard/Title";
import { UserVideos } from "./components/UserVideos";
import { LoadingVideos } from "./components/LoadingVideos";
import { useUserGeneratedVideos } from "./queries/useUserGeneratedVideos";

export const UserGeneratedVideos = () => {
  const {
    data: userVideos,
    isLoading: userVideosLoading,
    error: userVideosError,
  } = useUserGeneratedVideos();

  return (
    <div>
      <Title title="My Generated Videos" description="Videos you've created." />
      {userVideosLoading && <LoadingVideos />}
      {userVideosError && <div>Error: {userVideosError.message}</div>}
      {userVideos && <UserVideos userVideos={userVideos?.videos || []} />}
    </div>
  );
};
