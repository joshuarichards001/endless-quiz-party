export const getButtonColor = (id: number): string => {
  switch (id) {
    case 0:
      return "bg-blue-500 hover:bg-blue-700 text-white";
    case 1:
      return "bg-green-500 hover:bg-green-700 text-white";
    case 2:
      return "bg-yellow-500 hover:bg-yellow-700 text-black";
    case 3:
      return "bg-red-500 hover:bg-red-700 text-white";
    default:
      return "bg-gray-500 hover:bg-gray-700 text-white";
  }
};
