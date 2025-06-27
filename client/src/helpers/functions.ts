export const getButtonColor = (
  id: number,
  correctAnswer: number | null,
): string => {
  if (correctAnswer !== null) {
    if (id === correctAnswer) {
      return "bg-green-500 text-white border-green-700";
    } else {
      return "bg-red-500 text-white border-red-700 opacity-75";
    }
  }

  switch (id) {
    case 0:
      return "bg-blue-500 text-white border-blue-700";
    case 1:
      return "bg-green-500 text-white border-green-700";
    case 2:
      return "bg-yellow-500 text-black border-yellow-700";
    case 3:
      return "bg-red-500 text-white border-red-700";
    default:
      return "bg-gray-500 text-white border-gray-700";
  }
};
