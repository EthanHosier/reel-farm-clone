interface TitleProps {
  title: string;
  description: string;
}

export const Title = ({ title, description }: TitleProps) => {
  return (
    <div>
      <h2 className="text-xl font-semibold">{title}</h2>
      <p className="text-gray-600 mb-6 text-sm font-light">{description}</p>
    </div>
  );
};
