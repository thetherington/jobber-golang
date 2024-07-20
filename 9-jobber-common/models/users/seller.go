package users

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/thetherington/jobber-common/models/review"
	pbreview "github.com/thetherington/jobber-common/protogen/go/review"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Seller struct {
	Id               string                  `json:"_id"                        bson:"_id,omitempty"                `
	ProfilePublicId  string                  `json:"profilePublicId,omitempty"  bson:"profilePublicId,omitempty"    `
	FullName         string                  `json:"fullName"                   bson:"fullName"                     validate:"required"                        errmsg:"Full name required"`
	Username         string                  `json:"username,omitempty"         bson:"username,omitempty"           `
	Email            string                  `json:"email,omitempty"            bson:"email,omitepty"               validate:"required,email"                  errmsg:"Please enter a valid email address"`
	ProfilePicture   string                  `json:"profilePicture"             bson:"profilePicture"               validate:"required"                        errmsg:"Please add a profile picture"`
	Description      string                  `json:"description"                bson:"description"                  validate:"required"                        errmsg:"Please add a seller description"`
	Country          string                  `json:"country"                    bson:"country"                      validate:"required"                        errmsg:"Please select a country"`
	Oneliner         string                  `json:"oneliner"                   bson:"oneliner"                     validate:"required"                        errmsg:"Please add your oneliner"`
	Skills           []string                `json:"skills"                     bson:"skills"                       validate:"required,min=1"                  errmsg:"Please add at least one skill"`
	RatingsCount     int32                   `json:"ratingsCount"               bson:"ratingsCount"                 `
	RatingSum        int32                   `json:"ratingSum"                  bson:"ratingSum"                    `
	RatingCategories review.RatingCategories `json:"ratingCategories"           bson:"ratingCategories"             `
	Languages        []Language              `json:"languages"                  bson:"languages"                    validate:"required,dive"                   errmsg:"Please add at least one language"`
	ResponseTime     int32                   `json:"responseTime"               bson:"responseTime"                 validate:"required"                        errmsg:"Please add a response time"`
	RecentDelivery   *time.Time              `json:"recentDelivery"             bson:"recentDelivery"`
	Experience       []Experience            `json:"experience"                 bson:"experience"                   validate:"required,min=1,dive,required"    errmsg:"Please add at least one work experience"`
	Education        []Education             `json:"education"                  bson:"education"                    validate:"required,min=1,dive,required"    errmsg:"Please add at least one education"`
	SocialLinks      []string                `json:"socialLinks"                bson:"socialLinks"                  validate:"omitempty,min=0,dive,required"   errmsg:"Cannot have empty social link value"`
	Certificates     []Certificate           `json:"certificates"               bson:"certificates"                 validate:"omitempty,min=0,dive,required"   errmsg:"Certificate invalid"`
	OngoingJobs      int32                   `json:"ongoingJobs"                bson:"ongoingJobs"                  `
	CompletedJobs    int32                   `json:"completedJobs"              bson:"completedJobs"                `
	CancelledJobs    int32                   `json:"cancelledJobs"              bson:"cancelledJobs"                `
	TotalEarnings    float32                 `json:"totalEarnings"              bson:"totalEarnings,truncate"       `
	TotalGigs        int32                   `json:"totalGigs"                  bson:"totalGigs"                    `
	CreatedAt        time.Time               `json:"createdAt"                  bson:"createdAt"                    `
	UpdatedAt        time.Time               `json:"updatedAt"                  bson:"updatedAt"                    `
}

func (s *Seller) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[Seller](*s, validate)
}

type Language struct {
	Id       string `json:"_id"      bson:"_id,omitempty"`
	Language string `json:"language" bson:"language"   validate:"required" errmsg:"Invalid"`
	Level    string `json:"level"    bson:"level"      validate:"required" errmsg:"Invalid"`
}

type Experience struct {
	Id                   string `json:"_id"                  bson:"_id,omitempty"`
	Company              string `json:"company"              bson:"company"                validate:"required" errmsg:"Invalid"`
	Title                string `json:"title"                bson:"title"                  validate:"required" errmsg:"Invalid"`
	StartDate            string `json:"startDate"            bson:"startDate"              validate:"required" errmsg:"Invalid"`
	EndDate              string `json:"endDate"              bson:"endDate"                validate:"required" errmsg:"Invalid"`
	Description          string `json:"description"          bson:"description"            validate:"required" errmsg:"Invalid"`
	CurrentlyWorkingHere bool   `json:"currentlyWorkingHere" bson:"currentlyWorkingHere"`
}

type Education struct {
	Id         string `json:"_id"         bson:"_id,omitempty"`
	Country    string `json:"country"     bson:"country"    validate:"required" errmsg:"Invalid"`
	University string `json:"university"  bson:"university" validate:"required" errmsg:"Invalid"`
	Title      string `json:"title"       bson:"title"      validate:"required" errmsg:"Invalid"`
	Major      string `json:"major"       bson:"major"      validate:"required" errmsg:"Invalid"`
	Year       string `json:"year"        bson:"year"       validate:"required" errmsg:"Invalid"`
}

type Certificate struct {
	Id   string `json:"_id"    bson:"_id,omitempty"`
	Name string `json:"name"   bson:"name"    validate:"required" errmsg:"Invalid"`
	From string `json:"from"   bson:"from"    validate:"required" errmsg:"Invalid"`
	Year string `json:"year"   bson:"year"    validate:"required" errmsg:"Invalid"`
}

type SellerResponse struct {
	Message string  `json:"message"`
	Seller  *Seller `json:"seller"`
}

type SellersResponse struct {
	Message string    `json:"message"`
	Sellers []*Seller `json:"sellers"`
}

func createSetObjectID(id string) string {
	if id == "" {
		return primitive.NewObjectID().Hex()
	}

	return id
}

func CreateSellerFromReqPayload(s *pb.CreateUpdateSellerPayload) *Seller {
	seller := &Seller{
		ProfilePublicId: s.ProfilePublicId,
		FullName:        s.FullName,
		Username:        "",
		Email:           s.Email,
		ProfilePicture:  s.ProfilePicture,
		Description:     s.Description,
		Country:         s.Country,
		Oneliner:        s.Oneliner,
		ResponseTime:    s.ResponseTime,
		RatingsCount:    0,
		RatingSum:       0,
		OngoingJobs:     0,
		CompletedJobs:   0,
		CancelledJobs:   0,
		TotalEarnings:   0,
		TotalGigs:       0,
		RecentDelivery:  nil,
		SocialLinks:     make([]string, 0),
		Skills:          make([]string, 0),
		Languages:       make([]Language, 0),
		Experience:      make([]Experience, 0),
		Education:       make([]Education, 0),
		Certificates:    make([]Certificate, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		RatingCategories: review.RatingCategories{
			One:   review.RatingCategoryItem{},
			Two:   review.RatingCategoryItem{},
			Three: review.RatingCategoryItem{},
			Four:  review.RatingCategoryItem{},
			Five:  review.RatingCategoryItem{},
		},
	}

	seller.Skills = append(seller.Skills, s.Skills...)
	seller.SocialLinks = append(seller.SocialLinks, s.SocialLinks...)

	for _, v := range s.Languages {
		seller.Languages = append(seller.Languages, Language{
			Id:       createSetObjectID(*v.Id),
			Language: v.Language,
			Level:    v.Level,
		})
	}

	for _, v := range s.Experience {
		seller.Experience = append(seller.Experience, Experience{
			Id:                   createSetObjectID(*v.Id),
			Company:              v.Company,
			Title:                v.Title,
			StartDate:            v.StartDate,
			EndDate:              v.EndDate,
			Description:          v.Description,
			CurrentlyWorkingHere: v.CurrentlyWorkingHere,
		})
	}

	for _, v := range s.Education {
		seller.Education = append(seller.Education, Education{
			Id:         createSetObjectID(*v.Id),
			Country:    v.Country,
			University: v.University,
			Title:      v.Title,
			Major:      v.Major,
			Year:       v.Year,
		})
	}

	for _, v := range s.Certificates {
		seller.Certificates = append(seller.Certificates, Certificate{
			Id:   createSetObjectID(*v.Id),
			Name: v.Name,
			From: v.From,
			Year: v.Year,
		})
	}

	return seller
}

func CreateReqPayload(s *Seller) *pb.CreateUpdateSellerPayload {
	payload := &pb.CreateUpdateSellerPayload{
		Email:           s.Email,
		ProfilePublicId: s.ProfilePublicId,
		ProfilePicture:  s.ProfilePicture,
		FullName:        s.FullName,
		Description:     s.Description,
		Country:         s.Country,
		Oneliner:        s.Oneliner,
		ResponseTime:    s.ResponseTime,
		Skills:          make([]string, 0),
		SocialLinks:     make([]string, 0),
		Languages:       make([]*pb.Language, 0),
		Experience:      make([]*pb.Experience, 0),
		Education:       make([]*pb.Education, 0),
		Certificates:    make([]*pb.Certificate, 0),
	}

	payload.Skills = append(payload.Skills, s.Skills...)
	payload.SocialLinks = append(payload.SocialLinks, s.SocialLinks...)

	for _, v := range s.Languages {
		payload.Languages = append(payload.Languages, &pb.Language{
			Id:       &v.Id,
			Language: v.Language,
			Level:    v.Level,
		})
	}

	for _, v := range s.Experience {
		payload.Experience = append(payload.Experience, &pb.Experience{
			Id:                   &v.Id,
			Company:              v.Company,
			Title:                v.Title,
			StartDate:            v.StartDate,
			EndDate:              v.EndDate,
			Description:          v.Description,
			CurrentlyWorkingHere: v.CurrentlyWorkingHere,
		})
	}

	for _, v := range s.Education {
		payload.Education = append(payload.Education, &pb.Education{
			Id:         &v.Id,
			Country:    v.Country,
			University: v.University,
			Title:      v.Title,
			Major:      v.Major,
			Year:       v.Year,
		})
	}

	for _, v := range s.Certificates {
		payload.Certificates = append(payload.Certificates, &pb.Certificate{
			Id:   &v.Id,
			Name: v.Name,
			From: v.From,
			Year: v.Year,
		})
	}

	return payload
}

func CreatePayloadFromSeller(s *Seller) *pb.SellerPayload {
	payload := &pb.SellerPayload{
		Id:              s.Id,
		Username:        s.Username,
		Email:           s.Email,
		ProfilePublicId: s.ProfilePublicId,
		ProfilePicture:  s.ProfilePicture,
		FullName:        s.FullName,
		Description:     s.Description,
		Country:         s.Country,
		Oneliner:        s.Oneliner,
		ResponseTime:    s.ResponseTime,
		RatingsCount:    s.RatingsCount,
		RatingsSum:      s.RatingSum,
		OngoingJobs:     s.OngoingJobs,
		CompletedJobs:   s.CompletedJobs,
		CancelledJobs:   s.CancelledJobs,
		TotalEarnings:   s.TotalEarnings,
		TotalGigs:       s.TotalGigs,
		RecentDelivery:  utils.ToDateTimeOrNil(s.RecentDelivery),
		Skills:          make([]string, 0),
		SocialLinks:     make([]string, 0),
		Languages:       make([]*pb.Language, 0),
		Experience:      make([]*pb.Experience, 0),
		Education:       make([]*pb.Education, 0),
		Certificates:    make([]*pb.Certificate, 0),
		CreatedAt:       utils.ToDateTime(&s.CreatedAt),
		UpdatedAt:       utils.ToDateTime(&s.CreatedAt),
		RatingCategories: &pbreview.RatingCategories{
			One:   &pbreview.RatingCategoryItem{Value: s.RatingCategories.One.Value, Count: s.RatingCategories.One.Count},
			Two:   &pbreview.RatingCategoryItem{Value: s.RatingCategories.Two.Value, Count: s.RatingCategories.Two.Count},
			Three: &pbreview.RatingCategoryItem{Value: s.RatingCategories.Three.Value, Count: s.RatingCategories.Three.Count},
			Four:  &pbreview.RatingCategoryItem{Value: s.RatingCategories.Four.Value, Count: s.RatingCategories.Four.Count},
			Five:  &pbreview.RatingCategoryItem{Value: s.RatingCategories.Five.Value, Count: s.RatingCategories.Five.Count},
		},
	}

	payload.Skills = append(payload.Skills, s.Skills...)
	payload.SocialLinks = append(payload.SocialLinks, s.SocialLinks...)

	for _, v := range s.Languages {
		payload.Languages = append(payload.Languages, &pb.Language{
			Id:       &v.Id,
			Language: v.Language,
			Level:    v.Level,
		})
	}

	for _, v := range s.Experience {
		payload.Experience = append(payload.Experience, &pb.Experience{
			Id:                   &v.Id,
			Company:              v.Company,
			Title:                v.Title,
			StartDate:            v.StartDate,
			EndDate:              v.EndDate,
			Description:          v.Description,
			CurrentlyWorkingHere: v.CurrentlyWorkingHere,
		})
	}

	for _, v := range s.Education {
		payload.Education = append(payload.Education, &pb.Education{
			Id:         &v.Id,
			Country:    v.Country,
			University: v.University,
			Title:      v.Title,
			Major:      v.Major,
			Year:       v.Year,
		})
	}

	for _, v := range s.Certificates {
		payload.Certificates = append(payload.Certificates, &pb.Certificate{
			Id:   &v.Id,
			Name: v.Name,
			From: v.From,
			Year: v.Year,
		})
	}

	return payload
}

func CreateSellerFromPayload(s *pb.SellerPayload) *Seller {
	seller := &Seller{
		Id:              s.Id,
		ProfilePublicId: s.ProfilePublicId,
		FullName:        s.FullName,
		Username:        s.Username,
		Email:           s.Email,
		ProfilePicture:  s.ProfilePicture,
		Description:     s.Description,
		Country:         s.Country,
		Oneliner:        s.Oneliner,
		RatingsCount:    s.RatingsCount,
		RatingSum:       s.RatingsSum,
		ResponseTime:    s.ResponseTime,
		SocialLinks:     make([]string, 0),
		Skills:          make([]string, 0),
		Languages:       make([]Language, 0),
		Experience:      make([]Experience, 0),
		Education:       make([]Education, 0),
		Certificates:    make([]Certificate, 0),
		OngoingJobs:     s.OngoingJobs,
		CompletedJobs:   s.CompletedJobs,
		CancelledJobs:   s.CancelledJobs,
		TotalEarnings:   s.TotalEarnings,
		TotalGigs:       s.TotalGigs,
		RecentDelivery:  utils.ToTimeOrNil(s.RecentDelivery),
		CreatedAt:       *utils.ToTime(s.CreatedAt),
		UpdatedAt:       *utils.ToTime(s.UpdatedAt),
		RatingCategories: review.RatingCategories{
			One:   review.RatingCategoryItem{Value: s.RatingCategories.One.Value, Count: s.RatingCategories.One.Count},
			Two:   review.RatingCategoryItem{Value: s.RatingCategories.Two.Value, Count: s.RatingCategories.Two.Count},
			Three: review.RatingCategoryItem{Value: s.RatingCategories.Three.Value, Count: s.RatingCategories.Three.Count},
			Four:  review.RatingCategoryItem{Value: s.RatingCategories.Four.Value, Count: s.RatingCategories.Four.Count},
			Five:  review.RatingCategoryItem{Value: s.RatingCategories.Five.Value, Count: s.RatingCategories.Five.Count},
		},
	}

	seller.Skills = append(seller.Skills, s.Skills...)
	seller.SocialLinks = append(seller.SocialLinks, s.SocialLinks...)

	for _, v := range s.Languages {
		seller.Languages = append(seller.Languages, Language{
			Id:       *v.Id,
			Language: v.Language,
			Level:    v.Level,
		})
	}

	for _, v := range s.Experience {
		seller.Experience = append(seller.Experience, Experience{
			Id:                   *v.Id,
			Company:              v.Company,
			Title:                v.Title,
			StartDate:            v.StartDate,
			EndDate:              v.EndDate,
			Description:          v.Description,
			CurrentlyWorkingHere: v.CurrentlyWorkingHere,
		})
	}

	for _, v := range s.Education {
		seller.Education = append(seller.Education, Education{
			Id:         *v.Id,
			Country:    v.Country,
			University: v.University,
			Title:      v.Title,
			Major:      v.Major,
			Year:       v.Year,
		})
	}

	for _, v := range s.Certificates {
		seller.Certificates = append(seller.Certificates, Certificate{
			Id:   *v.Id,
			Name: v.Name,
			From: v.From,
			Year: v.Year,
		})
	}

	return seller
}
