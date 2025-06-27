export const getButtonColor = (
  id: number,
  correctAnswer: number | null,
): string => {
  if (correctAnswer !== null) {
    if (id === correctAnswer) {
      return "bg-gradient-to-br from-emerald-400 to-emerald-600 text-white border-emerald-700 shadow-lg transform scale-105";
    } else {
      return "bg-gradient-to-br from-gray-400 to-gray-600 text-white border-gray-700 opacity-60";
    }
  }

  switch (id) {
    case 0:
      return "bg-gradient-to-br from-blue-400 to-blue-600 text-white border-blue-700";
    case 1:
      return "bg-gradient-to-br from-emerald-400 to-emerald-600 text-white border-emerald-700";
    case 2:
      return "bg-gradient-to-br from-amber-400 to-orange-500 text-white border-orange-600";
    case 3:
      return "bg-gradient-to-br from-red-400 to-red-600 text-white border-red-700";
    default:
      return "bg-gradient-to-br from-gray-400 to-gray-600 text-white border-gray-700";
  }
};
