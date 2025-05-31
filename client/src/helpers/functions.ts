export const getButtonColor = (id: number): string => {
  switch (id) {
    case 0:
      return "bg-blue-500 hover:active:bg-blue-600 text-white border-blue-700";
    case 1:
      return "bg-green-500 hover:active:bg-green-600 text-white border-green-700";
    case 2:
      return "bg-yellow-500 hover:active:bg-yellow-600 text-black border-yellow-700";
    case 3:
      return "bg-red-500 hover:active:bg-red-600 text-white border-red-700";
    default:
      return "bg-gray-500 hover:active:bg-gray-600 text-white border-gray-700";
  }
};
